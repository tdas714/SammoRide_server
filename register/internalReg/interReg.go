package internalReg

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"sammoRide/ca"
)

func RegisterInter(RootCert *x509.Certificate, RootKey *ecdsa.PrivateKey, country, orgName, province, ipAddr string) {
	_, dcaBlock, priv := ca.GenDCA(RootCert, RootKey, country, orgName, ipAddr, province, int64(1))
	_ = os.Mkdir("interCerts", 0755)
	certOut, err := os.Create("peerCerts/interCa.crt")
	if err != nil {
		log.Panic("Pem file Creation failed: ", err)
	}
	//Public Key
	pem.Encode(certOut, dcaBlock)
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("interCerts/interCa.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")
}
