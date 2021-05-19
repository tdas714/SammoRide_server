package network

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

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
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
	}, nil
}

func SendData(addr, ca, crt, key string) {
	// addr := *connect
	if !strings.Contains(addr, ":") {
		addr += ":443"
	}

	config, err := createClientConfig(ca, crt, key)
	if err != nil {
		log.Panic("config failed: %s", err.Error())
	}

	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
	}
	defer conn.Close()

	log.Printf("connect to %s succeed", addr)
	// printConnState(conn)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		io.Copy(conn, os.Stdin)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		io.Copy(os.Stdout, conn)
		wg.Done()
	}()
	wg.Wait()
}
