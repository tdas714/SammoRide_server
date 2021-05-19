package internalReg

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"sammoRide/ca"
	"strings"
)

func RegisterInter(RootCert *x509.Certificate, RootKey *ecdsa.PrivateKey, country, orgName, province, ipAddr string) {
	dir := strings.ReplaceAll(orgName, ".", "/")
	_, dcaBlock, priv := ca.GenDCA(RootCert, RootKey, country, orgName, ipAddr, province, int64(1))
	_ = os.MkdirAll(fmt.Sprintf("interCerts/%s", dir), 0755)
	fmt.Println(dir)
	certOut, err := os.Create("interCerts/" + dir + "/interCa.crt")
	if err != nil {
		log.Panic("Pem file Creation failed: ", err)
	}
	//Public Key
	pem.Encode(certOut, dcaBlock)
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("interCerts/"+dir+"/interCa.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")
}
