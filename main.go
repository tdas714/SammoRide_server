package main

import (
	"flag"
	"fmt"
	"sammoRide/network"
	"sammoRide/ut/client"
)

func main() {
	fmt.Println("This is a test")
	mode := flag.String("mode", "Mode", "Mode of the application")
	flag.Parse()

	fmt.Print(*mode)
	// cert, err := ioutil.ReadFile("rootCerts/rootCa.crt")
	// key, err := ioutil.ReadFile("rootCerts/rootCa.key")
	// ut.CheckErr(err, "main")

	// internalReg.RegisterInter(ut.LoadCertificate(cert), ut.LoadPrivateKey(key), "India", "orderer", "west", "127.0.0.1", "kolkata", "100025")

	if *mode == "s" {
		go network.StartEnrollServer("interCerts/orderer")
		network.StartOrederServer("localhost", "rootCerts/rootCa.crt",
			"interCerts/orderer/interCa.crt",
			"interCerts/orderer/interCa.key")
	} else {
		// client.SendEnrollRequest("India", "Tapas.Das", "west Bengal", "kolkata", "700028")
		client.SendData("rootCerts/rootCa.crt", "PeerCerts/Cert.crt", "PeerCerts/Cert.key")
	}

}
