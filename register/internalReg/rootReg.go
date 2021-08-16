package internalReg

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"sammoRide/ca"
	"sammoRide/ut"
)

func RegisterRoot(mode string) (*x509.Certificate, *ecdsa.PrivateKey, []byte) {
	certfilepath := fmt.Sprintf("%sCerts/%sCa.crt", mode, mode)
	caBytes, priv := ca.GenCARoot()
	_ = os.Mkdir(fmt.Sprintf("%sCerts", mode), 0700)
	certOut, err := os.Create(certfilepath)
	ut.CheckErr(err, "RegisterRoot/certout")

	//Public Key
	pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	certOut.Close()
	log.Print("written cert.pem\n")

	privfilepath := fmt.Sprintf("%sCerts/%sCa.key", mode, mode)

	// Private key
	keyOut, err := os.OpenFile(privfilepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")

	// // Load CA
	catls, err := tls.LoadX509KeyPair(certfilepath, privfilepath)
	ut.CheckErr(err, "RegisterRoot/catls")
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	ut.CheckErr(err, "RegisterRoot/ca")

	pPem, err := os.ReadFile(privfilepath)
	ut.CheckErr(err, "REgisterRoot/pPem")
	p := ut.LoadPrivateKey(pPem)

	// RegisterInter(ca, p, "India", fmt.Sprintf("sammoride.orderer.%d.com", 1), "West Bengal", ut.GetIP(), "kolkata", "700028")
	return ca, p, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
}
