package ut

import (
	"crypto/x509"
	"encoding/pem"
	"log"
)

func VerifyOrderer(rootCa, ordererCa []byte) {

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootCa))
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode([]byte(ordererCa))
	if block == nil {
		panic("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	CheckErr(err, "VerifyOrderer/ParseCert")

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		panic("failed to verify certificate: " + err.Error())
	}
	log.Print("Orderer Verified")
}

func VerifyPeer(rootCa, ordererCa, peerCa []byte) {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootCa))
	if !ok {
		panic("failed to parse root certificate")
	}

	inters := x509.NewCertPool()
	ok = inters.AppendCertsFromPEM([]byte(ordererCa))
	if !ok {
		panic("failed to parse root certificate")
	}

	cert := LoadCertificate(peerCa)
	// CheckErr(err, "VerifyOrderer/ParseCert")

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inters,
	}

	if _, err := cert.Verify(opts); err != nil {
		panic("failed to verify certificate: " + err.Error())
	}
	log.Print("Peer Verified")

}
