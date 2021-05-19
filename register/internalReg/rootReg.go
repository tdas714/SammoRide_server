package internalReg

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sammoRide/ca"
	"sammoRide/ut"
)

func RegisterRoot() {
	_, rBlock, priv := ca.GenCARoot("India", "SammoRide", ut.GetIP(), int64(1))
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

	f, err := ioutil.ReadFile("rootCerts/rootCa.key")
	if err != nil {
		fmt.Println(err)
	}
	privateKey := ut.LoadPrivateKey(f)
	fmt.Println(privateKey)
	f, err = ioutil.ReadFile("rootCerts/rootCa.crt")
	cert := ut.LoadCertificate(f)

	RegisterInter(cert, privateKey, "India", fmt.Sprintf("sammoride.orderer.%d.com", 1), "West Bengal", ut.GetIP())
}
