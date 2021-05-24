package internalReg

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"sammoRide/ca"
	"sammoRide/ut"
	"strings"
)

func RegisterInter(RootCert *x509.Certificate, RootKey *ecdsa.PrivateKey, country, orgName, province, ipAddr, city, postalCode string) {
	dir := strings.ReplaceAll(orgName, ".", "/")

	sha := sha1.New()
	sha.Write([]byte("This will be ordrerer Req")) //<-- change that later

	ordererCaBytes, priv := ca.GenDCA(RootCert, RootKey, country, orgName, ipAddr, province, city, postalCode, int64(1), sha.Sum(nil))

	_ = os.MkdirAll(fmt.Sprintf("interCerts/%s", dir), 0700)

	certOut, err := os.Create("interCerts/" + dir + "/interCa.crt")
	if err != nil {
		log.Panic("Pem file Creation failed: ", err)
	}

	//Public Key
	pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ordererCaBytes,
	})
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("interCerts/"+dir+"/interCa.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	ut.CheckErr(err, "RegisterInter/privBytes")
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")
}
