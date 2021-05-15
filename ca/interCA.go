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

func GenDCA(RootCert *x509.Certificate, RootKey *ecdsa.PrivateKey, country, orgName, ipAddr, province string, serNum int64) (*x509.Certificate, *pem.Block, *ecdsa.PrivateKey) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	var DCATemplate = x509.Certificate{
		SerialNumber: big.NewInt(serNum),
		Subject: pkix.Name{
			Country:            []string{country},
			Organization:       []string{orgName},
			Province:           []string{province},
			OrganizationalUnit: []string{"intermediate"},
			CommonName:         "DCA",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(25, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        false,
		MaxPathLen:            1,
		IPAddresses:           []net.IP{net.ParseIP(ipAddr)},
	}
	DCACert, DCABlock := genCert(&DCATemplate, RootCert, &priv.PublicKey, RootKey)
	return DCACert, DCABlock, priv
}
