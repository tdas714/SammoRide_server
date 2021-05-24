package ut

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	ENROLL_REQ = "EnrollRequest"
	ENROLL_RES = "EnrollResponce"
	SERIAL_LOG = "rootCerts/Serial.log"
)

type PeerEnrollDataRequest struct {
	Country    string
	Name       string
	Province   string
	IpAddr     string
	City       string
	PostalCode string
}

type PeerEnrollDataResponse struct {
	Header     string
	IpAddr     string
	PeerCert   []byte
	PrivateKey []byte
	SenderCert []byte
	RootCert   []byte
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	l := make([]string, 0)
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				l = append(l, ipnet.IP.String())
			}
		}
	}
	return l[0]
}

func LoadPrivateKey(f []byte) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(string(f)))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}

func LoadCertificate(f []byte) *x509.Certificate {
	block, _ := pem.Decode(f)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}
	return cert
}

func CheckErr(err error, origin string) {
	if err != nil {
		log.Fatalf("%s - %s", origin, err)
	}
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
