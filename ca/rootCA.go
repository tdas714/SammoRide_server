package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sammoRide/ut"
	"time"
)

func GenCARoot() ([]byte, *ecdsa.PrivateKey) {
	// if _, err := os.Stat("someFile"); err == nil {
	// 	//read PEM and cert from file
	// }
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"Sammo Ride pvt. ltd."},
			Country:       []string{"India"},
			StreetAddress: []string{""},
			PostalCode:    []string{"700028"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1000, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		// IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ut.CheckErr(err, "GenCARoot/priv")

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &priv.PublicKey, priv)
	ut.CheckErr(err, "GenCARoot/caBytes")

	return caBytes, priv
}
