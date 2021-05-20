package network

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sammoRide/ut"

	"github.com/dgraph-io/badger/v3"
)

func createOrdererServerConfig(caPath, crtPath, keyPath string) (*tls.Config, error) {
	caCertPEM, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}, nil
}

func StartOrederServer(ipAddr, caPath, crtPath, keyPath string) {
	config, err := createOrdererServerConfig(caPath, crtPath, keyPath)
	if err != nil {
		log.Fatal("config failed: %s", err.Error())
	}

	ln, err := tls.Listen("tcp", ipAddr, config)
	if err != nil {
		log.Fatal("listen failed: %s", err.Error())
	}

	log.Printf("listen on %s", ipAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept failed: %s", err.Error())
			break
		}
		log.Printf("connection open: %s", conn.RemoteAddr())
		// printConnState(conn.(*tls.Conn))

		go func(c net.Conn) {
			wr, _ := io.Copy(c, c)
			c.Close()
			log.Printf("connection close: %s, written: %d", conn.RemoteAddr(), wr)
		}(conn)
	}
}

func StartEnrollServer(name string) {
	// var enrollReq *ut.PeerEnrollDataRequest
	// cert, err := ioutil.ReadFile(name + "/interCa.crt")
	// rcert, err := ioutil.ReadFile("rootCerts/rootCa.crt")
	// priv, err := ioutil.ReadFile(name + "/interCa.key")
	// ut.CheckErr(err)

	db, err := badger.Open(badger.DefaultOptions("database"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b := make([]byte, 1024)
	if _, err := os.Stat(ut.SERIAL_LOG); os.IsNotExist(err) {
		binary.BigEndian.PutUint64(b, 0)
		err = ioutil.WriteFile(ut.SERIAL_LOG, b, 0700)
		ut.CheckErr(err)
	}

	// serialNum, err := ioutil.ReadFile(ut.SERIAL_LOG)
	ut.CheckErr(err)
	path := "/post"
	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		// handleRequest(enrollReq, cert, priv, rcert, serialNum, db, w, r)
		var res *ut.PeerEnrollDataRequest

		json.NewDecoder(r.Body).Decode(&res)
		fmt.Println(res.Name)
	})
	http.ListenAndServe("localhost:8080", nil)

}

func handleRequest(enrollReq *ut.PeerEnrollDataRequest, cert, priv, rcert, serialNum []byte,
	db *badger.DB, r *http.Request, rw http.ResponseWriter) {
	// enrollReq = r.Body.
	// _, pBlock, pPriv := ca.GenServerCert(ut.LoadCertificate(cert),
	// 	ut.LoadPrivateKey(priv),
	// 	enrollReq.IpAddr,
	// 	enrollReq.Country,
	// 	enrollReq.Name,
	// 	enrollReq.Province,
	// 	"Peer",
	// 	int64(binary.BigEndian.Uint64(serialNum)+1))

	// enrollRes := ut.PeerEnrollDataResponse{Header: ut.ENROLL_RES,
	// 	IpAddr:        enrollReq.IpAddr,
	// 	PeerCertBlock: *pBlock,
	// 	PrivateKey:    *pPriv,
	// 	SenderCert:    *ut.LoadCertificate(cert),
	// 	RootCert:      *ut.LoadCertificate(rcert)}

	// // j, err := json.Marshal(enrollRes)
	// if err != nil {
	// 	return
	// }
	// t := time.Now()
	// myTime := t.Format(time.RFC3339) + "\n"
	// cServer.Write([]byte(myTime))
	// fmt.Println(enrollReq)

	//
	// b := make([]byte, 1024)
	// cServer.Write(append(j, byte('\n')))
	// err = db.Update(func(txn *badger.Txn) error {
	// 	binary.BigEndian.PutUint64(b, binary.BigEndian.Uint64(serialNum)+1)
	// 	interByte, e := ut.GetBytes(enrollReq)
	// 	txn.Set([]byte(b), interByte)
	// 	ut.CheckErr(e)
	// 	return nil
	// })
	// ut.CheckErr(err)

}
