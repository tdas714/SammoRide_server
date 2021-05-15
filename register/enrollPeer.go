package register

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"sammoRide/ca"
)

func RegisterPeer(country, name, province, ipAddr string, dcaCert *x509.Certificate, dcaKey *ecdsa.PrivateKey) {
	_, peerBlock, priv := ca.GenServerCert(dcaCert, dcaKey, ipAddr, country, name, province, "peer", int64(1))

	_ = os.Mkdir("peerCerts", 0755)
	certOut, err := os.Create("peerCerts/peerCa.crt")
	if err != nil {
		log.Panic("Pem file Creation failed: ", err)
	}
	//Public Key
	pem.Encode(certOut, peerBlock)
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("PeerCerts/peerCa.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privbyte, err := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privbyte})
	keyOut.Close()
	log.Print("written key.pem\n")
}
