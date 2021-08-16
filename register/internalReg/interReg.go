package internalReg

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sammoRide/ca"
	"sammoRide/ut"
)

func RegisterInter(RootCert *x509.Certificate,
	RootKey *ecdsa.PrivateKey, country, orgName, province,
	ipAddr, city, postalCode, mode string) (*x509.Certificate, *ecdsa.PrivateKey, []byte) {
	certpath := fmt.Sprintf("%sCerts/%sCa.crt", mode, mode)
	keyPath := fmt.Sprintf("%sCerts/%sCa.key", mode, mode)
	sha := sha1.New()
	sha.Write([]byte("This will be ordrerer Req")) //<-- change that later

	ordererCaBytes, priv := ca.GenDCA(mode, RootCert, RootKey, country, orgName, ipAddr, province, city, postalCode, int64(1), sha.Sum(nil))

	_ = os.MkdirAll(fmt.Sprintf("%sCerts/", mode), 0700)

	certOut, err := os.Create(certpath)
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
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	ut.CheckErr(err, "RegisterInter/privBytes")
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")

	// catls, err := tls.LoadX509KeyPair(certpath, keyPath)
	// ut.CheckErr(err, "Registerinter/catls")
	// ca, err := x509.ParseCertificate(catls.Certificate[0])
	// ut.CheckErr(err, "Registerinter/ca")

	pPem, err := os.ReadFile(keyPath)
	ut.CheckErr(err, "REgisterinter/pPem")
	p := ut.LoadPrivateKey(pPem)

	cabytrs, err := ioutil.ReadFile(certpath)
	ut.CheckErr(err, "inuterReg/Loadcert")
	ca := ut.LoadCertificate(cabytrs)

	return ca, p, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ordererCaBytes,
	})
}
