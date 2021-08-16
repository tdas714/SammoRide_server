package network

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sammoRide/ca"
	"sammoRide/database"
	"sammoRide/struct/common"
	"sammoRide/struct/orderer"
	"sammoRide/ut"
	"strings"
)

func createOrdererServerConfig(caPath, rcaPath, crtPath, keyPath string) (*tls.Config, error) {
	caCertPEM, err := ioutil.ReadFile(caPath)
	ut.CheckErr(err, "createOrderconfig/caPem")
	RCertPem, err := ioutil.ReadFile(rcaPath)
	ut.CheckErr(err, "createOrdererConfig/certpem")

	roots, err := x509.SystemCertPool()
	ut.CheckErr(err, "create Order Config/roots")
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}
	ok = roots.AppendCertsFromPEM(RCertPem)
	ut.CheckErr(err, "createOrdererConfig/CertsPemAppend")

	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, err
	}
	//This Needs ATTENTION!!!!!
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

func StartOrederServer(ipAddr string, db *database.Database) {
	defer db.Close()
	// Set up a /hello resource handler
	http.HandleFunc("/hello", HelloHandler)
	http.HandleFunc("/Announcement/rider", func(rw http.ResponseWriter, r *http.Request) {

		RiderAHandler(rw, r, db)
	})

	http.HandleFunc("/TransactionCommitmentRequest", func(rw http.ResponseWriter, r *http.Request) {
		TransactionCommitmentHandle(rw, r, db)
	})

	http.HandleFunc("/Enroll/Orderer", func(rw http.ResponseWriter, r *http.Request) {
		OrdererEnrollHandler(rw, r, db)
	})

	http.HandleFunc("/Request/BlockSnapshot", func(rw http.ResponseWriter, r *http.Request) {
		SnapshotHandler(rw, r, db)
	})

	tlsConfig, err := createOrdererServerConfig(db.InterCaPath, db.RootCaPath, db.Certificatepath, db.KeyPath)
	ut.CheckErr(err, "StartOrederServer/config")

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ipAddr + db.Info.Port,
		TLSConfig: tlsConfig,
	}
	fmt.Println("Server Starting : ", ipAddr+db.Info.Port)
	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS(db.Certificatepath, db.KeyPath))
}

func StartEnrollServer(db *database.Database) {
	defer db.Close()
	// var enrollReq *ut.PeerEnrollDataRequest
	cert, err := ioutil.ReadFile(db.Certificatepath)
	rcert, err := ioutil.ReadFile(db.InterCaPath)
	priv, err := ioutil.ReadFile(db.KeyPath)
	ut.CheckErr(err, "StartEnrollServer")

	b := make([]byte, 1024)
	if _, err := os.Stat(ut.SERIAL_LOG); os.IsNotExist(err) {
		binary.BigEndian.PutUint64(b, 0)
		err = ioutil.WriteFile(strings.Split(db.UtilsPath, "/")[0]+ut.SERIAL_LOG, b, 0700)
		ut.CheckErr(err, "StartEnrollServer")
	}

	// fmt.Print(db.DB.IsClosed())

	serialNum, err := ioutil.ReadFile(strings.Split(db.UtilsPath, "/")[0] + ut.SERIAL_LOG)
	ut.CheckErr(err, "StartEnrollServer Serial")
	path := "/post"
	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		var res *orderer.EnrollDataRequest

		json.NewDecoder(r.Body).Decode(&res)
		handleRequest(res, cert, priv, rcert, serialNum, rw, db)
	})
	http.ListenAndServe(db.Info.IP+":8080", nil)

}

func handleRequest(enrollReq *orderer.EnrollDataRequest, cert, priv, rcert, serialNum []byte,
	rw http.ResponseWriter, db *database.Database) {

	fmt.Println("Enroll Request from ", enrollReq.IpAddr, enrollReq.ListingPort)

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
	)

	pPrivByte, err := x509.MarshalECPrivateKey(pPriv)
	ut.CheckErr(err, "handleReq/pPriveByte")

	enrollReq.PrivateKey = pPrivByte

	db.InsertNode(ut.GetBytes(enrollReq.Name+":"+enrollReq.Country+":"+enrollReq.Province+":"+enrollReq.City),
		enrollReq.Serialize()) //Later will be e-mail

	pCertPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: pCert,
	})

	enrollRes := orderer.EnrollDataResponse{Header: ut.ENROLL_RES,
		IpAddr:      enrollReq.IpAddr,
		PeerCert:    pCertPem,
		PrivateKey:  pPrivByte,
		SenderCert:  cert,
		RootCert:    rcert,
		PeerList:    db.GetRandomPeer(1),
		OrdererList: db.GetRandomOrderer(1),
	}

	js, err := json.Marshal(enrollRes)
	// fmt.Println(string(js))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
	fmt.Println(strings.Split(db.UtilsPath, "/")[0])

	err = ioutil.WriteFile(strings.Split(db.UtilsPath, "/")[0]+ut.SERIAL_LOG,
		ut.IntToByteArray(int64(binary.BigEndian.Uint64(serialNum)+1)), 0700)
	ut.CheckErr(err, "WriteSerialNumber")
}

func SendData(ca, crt, key, ipAddr, port,
	reqSubDomain string, data []byte) {

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
	// r, err := client.Get(fmt.Sprintf("https://%s/hello", addr))
	// CheckErr(err, "SendData/r")

	// =======POST
	url := fmt.Sprintf("https://%s:%s/%s", ipAddr,
		port, reqSubDomain)
	r, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err == nil {
		// Read the response body
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Print the response body to stdout
		fmt.Printf("%s\n", body)
	}
}

func SendEnrollRequest(ord *common.OrdererInfo, serverIp, serverPort string, paths *common.FilePath) []string {
	enrollReq := &orderer.EnrollDataRequest{Country: ord.Country, Name: ord.Name, Province: ord.Province, IpAddr: ord.Name,
		City: ord.City, PostalCode: ord.Postalcode, ListingPort: ord.Port}
	json_data, err := json.Marshal(enrollReq)
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf(fmt.Sprintf("http://%s:%s/Enroll/Orderer",
		serverIp, serverPort))

	resp, err := http.Post(path, "application/json", bytes.NewBuffer(json_data))
	ut.CheckErr(err, "SendEnrollRequest/Post")

	var res *orderer.EnrollDataResponse
	json.NewDecoder(resp.Body).Decode(&res)

	ut.VerifyPeer(ut.LoadCertificate(res.RootCert), ut.LoadCertificate(res.SenderCert), ut.LoadCertificate(res.PeerCert))
	ut.VerifyOrderer(ut.LoadCertificate(res.RootCert), ut.LoadCertificate(res.SenderCert))
	//

	interca := fmt.Sprintf("%s/interCa.crt", paths.CAsPath)
	rootca := fmt.Sprintf("%s/rootCa.crt", paths.CAsPath)

	_ = os.Mkdir(paths.CAsPath, 0700)
	_ = os.Mkdir(paths.CertificatePath, 0700)
	_ = os.Mkdir(paths.KeyPath, 0700)

	err = ioutil.WriteFile(interca, res.SenderCert, 0700)
	ut.CheckErr(err, "SendErollReq/interca")
	err = ioutil.WriteFile(rootca, res.RootCert, 0700)
	ut.CheckErr(err, "SendErollReq/rootca")

	blocks, _ := pem.Decode(res.PeerCert)
	if blocks == nil {
		log.Panic("Block is nil")
	}
	// Have to change Certificate Path in Yaml
	// Public key
	certOut, err := os.Create(paths.CertificatePath)

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: blocks.Bytes})
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private Key
	keyOut, err := os.OpenFile(paths.KeyPath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	ut.CheckErr(err, "SendEnrollRequest/keyOut")
	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: res.PrivateKey})
	ut.CheckErr(err, "SendEnrollRequest/penEncode")
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}
	log.Print("wrote key.pem\n")

	return res.PeerList
}

func StartService(YamlPath string, isGenesis bool) {

	InitFileStructure(YamlPath)

	var input *common.InputInfo
	odr, paths := input.Parse(YamlPath)

	interca := fmt.Sprintf("%s/interCa.crt", paths.CAsPath)
	rootca := fmt.Sprintf("%s/rootCa.crt", paths.CAsPath)

	ServerDatabase := database.NewDatabase(odr, paths.PeerDBPath, paths.OrdererfileDB, interca, rootca,
		paths.CertificatePath, paths.KeyPath, paths.ChainDB, paths.StatePath, isGenesis)

	go StartEnrollServer(ServerDatabase)
	StartOrederServer(odr.IP, ServerDatabase)
}

func EnrollOrderer(YamlPath string) {

	InitFileStructure(YamlPath)

	var input *common.InputInfo
	odr, paths := input.Parse(YamlPath)
	SendEnrollRequest(odr, "127.0.0.1", "8080", paths)

	interca := fmt.Sprintf("%s/interCa.crt", paths.CAsPath)
	rootca := fmt.Sprintf("%s/rootCa.crt", paths.CAsPath)

	ServerDatabase := database.NewDatabase(odr, paths.PeerDBPath, paths.OrdererfileDB,
		rootca, interca, paths.CertificatePath,
		paths.KeyPath, paths.ChainDB, paths.StatePath, true)

	ServerDatabase.Close()
}

func InitFileStructure(yamlfile string) {

	var input *common.InputInfo
	_, paths := input.Parse(yamlfile)
	dirCert := strings.Split(paths.CertificatePath, "/")[0]
	_ = os.Mkdir(dirCert, 0700)
	dirKey := strings.Split(paths.KeyPath, "/")[0]
	_ = os.Mkdir(dirKey, 0700)
	_ = os.Mkdir(paths.CAsPath, 0700)

	dir := strings.Split(paths.PeerDBPath, "/")[0]
	_ = os.Mkdir(dir, 0700)

	dir = strings.Split(paths.ChainDB, "/")[0]
	_ = os.Mkdir(dir, 0700)

	dir = strings.Split(paths.OrdererfileDB, "/")[0]
	_ = os.Mkdir(dir, 0700)

	dir = strings.Split(paths.StatePath, "/")[0]
	_ = os.Mkdir(dir, 0700)
	_ = os.Mkdir(paths.StatePath, 0700)
}

// "rootCerts/rootCa.crt"
// "interCerts/orderer/interCa.crt"
// "interCerts/orderer/interCa.key"
