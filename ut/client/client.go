package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sammoRide/ut"
)

func SendEnrollRequest(country, name, province, city, postC string) {
	enrollReq := &ut.PeerEnrollDataRequest{Country: country, Name: name, Province: province, IpAddr: "127.0.0.1",
		City: city, PostalCode: postC}
	json_data, err := json.Marshal(enrollReq)
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("http://localhost:8080/post")
	resp, err := http.Post(path, "application/json", bytes.NewBuffer(json_data))
	ut.CheckErr(err, "SendEnrollRequest/Post")

	var res *ut.PeerEnrollDataResponse
	json.NewDecoder(resp.Body).Decode(&res)

	ut.VerifyPeer(res.RootCert, res.SenderCert, res.PeerCert)
	ut.VerifyOrderer(res.RootCert, res.SenderCert)
	//
	_ = os.Mkdir("PeerCerts", 0700)

	blocks, _ := pem.Decode(res.PeerCert)

	// Public key
	certOut, err := os.Create("PeerCerts/Cert.crt")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: blocks.Bytes})
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private Key
	keyOut, err := os.OpenFile("PeerCerts/Cert.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	ut.CheckErr(err, "SendEnrollRequest/keyOut")
	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: res.PrivateKey})
	ut.CheckErr(err, "SendEnrollRequest/penEncode")
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}
	log.Print("wrote key.pem\n")
}

// ====================================
func createClientConfig(ca, crt, key string) (*tls.Config, error) {
	caCertPEM, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(crt, key)
	ut.CheckErr(err, "createClientConfig/cert")

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
	}, nil
}

func SendData(ca, crt, key string) {
	// addr := *connect
	// if !strings.Contains(addr, ":") {
	// 	addr += ":8443"
	// }

	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair(crt, key)
	ut.CheckErr(err, "SendData/cert")

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile(ca)
	ut.CheckErr(err, "SendData/cacert")

	caCertPool, err := x509.SystemCertPool()
	ut.CheckErr(err, "Client/CertPool")
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a HTTPS client and supply the created CA pool and certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		},
	}

	// Request /hello via the created HTTPS client over port 8443 via GET
	r, err := client.Get(fmt.Sprintf("https://%s:8443/hello", "127.0.0.1"))
	ut.CheckErr(err, "SendData/r")

	// Read the response body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("%s\n", body)
}
