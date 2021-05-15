package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

func genCert(template, parent *x509.Certificate, publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) (*x509.Certificate, *pem.Block) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
	if err != nil {
		panic("Failed to create certificate:" + err.Error())
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		panic("Failed to parse certificate:" + err.Error())
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	// certPEM := pem.EncodeToMemory(&b)

	return cert, &b
}

func GenServerCert(DCACert *x509.Certificate, DCAKey *ecdsa.PrivateKey, ipAddr, country, orgName, province, orgUnit string, serNum int64) (*x509.Certificate, *pem.Block, *ecdsa.PrivateKey) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	var ServerTemplate = x509.Certificate{
		SerialNumber: big.NewInt(serNum),
		Subject: pkix.Name{
			Country:            []string{country},
			Organization:       []string{orgName},
			Province:           []string{province},
			OrganizationalUnit: []string{orgUnit},
			CommonName:         "Peers",
		},
		NotBefore:      time.Now().Add(-10 * time.Second),
		NotAfter:       time.Now().AddDate(10, 0, 0),
		KeyUsage:       x509.KeyUsageCRLSign,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:           false,
		MaxPathLenZero: true,
		IPAddresses:    []net.IP{net.ParseIP(ipAddr)},
	}

	ServerCert, ServerBlock := genCert(&ServerTemplate, DCACert, &priv.PublicKey, DCAKey)
	return ServerCert, ServerBlock, priv

}
