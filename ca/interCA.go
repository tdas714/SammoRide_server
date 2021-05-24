package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

func GenDCA(rootCa *x509.Certificate, rootKey *ecdsa.PrivateKey, country, orgName, ipAddr, province, city, postalCode string, serNum int64, subKeyId []byte) ([]byte, *ecdsa.PrivateKey) {

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(serNum),
		Subject: pkix.Name{
			Organization:  []string{orgName},
			Country:       []string{country},
			Province:      []string{province},
			Locality:      []string{city},
			StreetAddress: []string{""},
			PostalCode:    []string{postalCode},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(25, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP(ipAddr)},
		SubjectKeyId:          subKeyId,
	}

	// extSubjectAltName := pkix.Extension{}
	// extSubjectAltName.Id = asn1.ObjectIdentifier{2, 5, 29, 17}
	// extSubjectAltName.Critical = true
	// extSubjectAltName.Value = []byte(`IP:` + ipAddr)
	// cert.ExtraExtensions = []pkix.Extension{extSubjectAltName}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	ordererCaBytes, err := x509.CreateCertificate(rand.Reader, cert, rootCa, &priv.PublicKey, rootKey)

	return ordererCaBytes, priv
}
