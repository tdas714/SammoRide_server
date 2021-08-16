package ut

import (
	"crypto/x509"
	"fmt"
)

func VerifyOrderer(rootCa, ordererCa *x509.Certificate) {

	roots := x509.NewCertPool()
	roots.AddCert(rootCa)
	opts := x509.VerifyOptions{
		Roots: roots,
	}

	if _, err := ordererCa.Verify(opts); err != nil {
		panic("failed to verify certificate: " + err.Error())
	}
	fmt.Println("Orderer verified")
}

func VerifyPeer(rootCa, ordererCa, peerCa *x509.Certificate) {
	roots := x509.NewCertPool()
	inter := x509.NewCertPool()
	roots.AddCert(rootCa)
	inter.AddCert(ordererCa)
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inter,
	}

	if _, err := peerCa.Verify(opts); err != nil {
		panic("failed to verify certificate: " + err.Error())
	}
	fmt.Println("Peer Verified")

}
