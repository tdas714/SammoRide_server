package network

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
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

func startOrederServer(ipAddr, caPath, crtPath, keyPath string) {
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
