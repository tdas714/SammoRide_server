package network

import (
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sammoRide/ca"
	"sammoRide/ut"

	"github.com/dgraph-io/badger/v3"
)

func createOrdererServerConfig(caPath, crtPath, keyPath string) (*tls.Config, error) {
	caCertPEM, err := ioutil.ReadFile(caPath)
	ut.CheckErr(err, "createOrderconfig/caPem")
	CertPem, err := ioutil.ReadFile(crtPath)
	ut.CheckErr(err, "createOrdererConfig/certpem")

	roots, err := x509.SystemCertPool()
	ut.CheckErr(err, "create Order Config/roots")
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}
	ok = roots.AppendCertsFromPEM(CertPem)
	ut.CheckErr(err, "createOrdererConfig/CertsPemAppend")

	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	io.WriteString(w, "Hello, world!\n")
}

func StartOrederServer(ipAddr, caPath, crtPath, keyPath string) {

	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	tlsConfig, err := createOrdererServerConfig(caPath, crtPath, keyPath)
	ut.CheckErr(err, "StartOrederServer/config")

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ipAddr + ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS(crtPath, keyPath))
}

func StartEnrollServer(name string) {
	// var enrollReq *ut.PeerEnrollDataRequest
	cert, err := ioutil.ReadFile(name + "/interCa.crt")
	rcert, err := ioutil.ReadFile("rootCerts/rootCa.crt")
	priv, err := ioutil.ReadFile(name + "/interCa.key")
	ut.CheckErr(err, "StartEnrollServer")

	db, err := badger.Open(badger.DefaultOptions("database"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b := make([]byte, 1024)
	if _, err := os.Stat(ut.SERIAL_LOG); os.IsNotExist(err) {
		binary.BigEndian.PutUint64(b, 0)
		err = ioutil.WriteFile(ut.SERIAL_LOG, b, 0700)
		ut.CheckErr(err, "StartEnrollServer")
	}

	serialNum, err := ioutil.ReadFile(ut.SERIAL_LOG)
	ut.CheckErr(err, "StartEnrollServer")
	path := "/post"
	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		var res *ut.PeerEnrollDataRequest

		json.NewDecoder(r.Body).Decode(&res)
		handleRequest(res, cert, priv, rcert, serialNum, db, rw)
	})
	http.ListenAndServe("localhost:8080", nil)

}

func handleRequest(enrollReq *ut.PeerEnrollDataRequest, cert, priv, rcert, serialNum []byte,
	db *badger.DB, rw http.ResponseWriter) {

	sha := sha1.New()
	js, err := json.Marshal(enrollReq)
	ut.CheckErr(err, "handleRequest/js")
	sha.Write(js)

	pCert, pPriv := ca.GenServerCert(ut.LoadCertificate(cert),
		ut.LoadPrivateKey(priv),
		enrollReq.IpAddr,
		enrollReq.Country,
		enrollReq.Name,
		enrollReq.Province,
		"Peer",
		enrollReq.City,
		enrollReq.PostalCode,
		int64(binary.BigEndian.Uint64(serialNum)+1),
		sha.Sum(nil),
	)

	pPrivByte, err := x509.MarshalECPrivateKey(pPriv)
	ut.CheckErr(err, "handleReq/pPriveByte")

	pCertPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: pCert,
	})

	_ = ut.LoadCertificate(pCertPem)

	enrollRes := ut.PeerEnrollDataResponse{Header: ut.ENROLL_RES,
		IpAddr:     enrollReq.IpAddr,
		PeerCert:   pCertPem,
		PrivateKey: pPrivByte,
		SenderCert: cert,
		RootCert:   rcert}

	//
	js, err = json.Marshal(enrollRes)
	// fmt.Println(string(js))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}
