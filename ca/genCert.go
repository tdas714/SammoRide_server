package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"net"
	"sammoRide/ut"
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

func GenServerCert(DCACert *x509.Certificate, DCAKey *ecdsa.PrivateKey, ipAddr, country, orgName, province, orgUnit, city, postalCode string, serNum int64, subKeyId []byte) ([]byte, *ecdsa.PrivateKey) {

	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(serNum),
		Subject: pkix.Name{
			Organization:       []string{orgName},
			Country:            []string{country},
			Province:           []string{province},
			Locality:           []string{city},
			StreetAddress:      []string{""},
			PostalCode:         []string{postalCode},
			OrganizationalUnit: []string{orgUnit},
			CommonName:         "IP:" + ut.GetIP(),
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: subKeyId,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		IPAddresses:  []net.IP{net.ParseIP(ipAddr)},
	}

	extSubjectAltName := pkix.Extension{}
	extSubjectAltName.Id = asn1.ObjectIdentifier{2, 5, 29, 17}
	extSubjectAltName.Critical = true
	extSubjectAltName.Value = []byte(`IP:` + ut.GetIP())
	cert.ExtraExtensions = []pkix.Extension{extSubjectAltName}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ut.CheckErr(err, "GenServerCert/privebyte")

	serverCert, err := x509.CreateCertificate(rand.Reader, cert, DCACert, &priv.PublicKey, DCAKey)
	ut.CheckErr(err, "GenServerCert/serverCert")

	return serverCert, priv

}
