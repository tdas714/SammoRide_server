package internalReg

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"sammoRide/ca"
)

func RegisterRoot() {
	_, rBlock, priv := ca.GenCARoot("India", "SammoRide", "127.0.0.1", int64(1))
	_ = os.Mkdir("rootCerts", 0700)
	certOut, err := os.Create("rootCerts/rootCa.crt")
	if err != nil {
		log.Panic("Pem file Creation failed: ", err)
	}
	//Public Key
	pem.Encode(certOut, rBlock)
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("rootCerts/rootCa.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")
}
