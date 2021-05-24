package internalReg

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"sammoRide/ca"
	"sammoRide/ut"
)

func RegisterRoot() {
	caBytes, priv := ca.GenCARoot()
	_ = os.Mkdir("rootCerts", 0600)
	certOut, err := os.Create("rootCerts/rootCa.crt")
	ut.CheckErr(err, "RegisterRoot/certout")

	//Public Key
	pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("rootCerts/rootCa.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")

	// Load CA
	catls, err := tls.LoadX509KeyPair("rootCerts/rootCa.crt", "rootCerts/rootCa.key")
	ut.CheckErr(err, "RegisterRoot/catls")
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	ut.CheckErr(err, "RegisterRoot/ca")

	pPem, err := os.ReadFile("rootCerts/rootCa.key")
	ut.CheckErr(err, "REgisterRoot/pPem")
	p := ut.LoadPrivateKey(pPem)

	RegisterInter(ca, p, "India", fmt.Sprintf("sammoride.orderer.%d.com", 1), "West Bengal", ut.GetIP(), "kolkata", "700028")

}
