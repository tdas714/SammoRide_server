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

func GenCARoot(country, orgName, ipAddr string, serNum int64) (*x509.Certificate, *pem.Block, *ecdsa.PrivateKey) {
	// if _, err := os.Stat("someFile"); err == nil {
	// 	//read PEM and cert from file
	// }
	var rootTemplate = x509.Certificate{
		SerialNumber: big.NewInt(serNum),
		Subject: pkix.Name{
			Country:      []string{country},
			Organization: []string{orgName},
			CommonName:   "Sammo Ride Root CA",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(25, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		IPAddresses:           []net.IP{net.ParseIP(ipAddr)},
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	rootCert, rootBlock := genCert(&rootTemplate, &rootTemplate, &priv.PublicKey, priv)
	return rootCert, rootBlock, priv
}
