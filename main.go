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
	// _, rootCertPEM, _ := ca.GenCARoot("India", "SammoRide", "127.0.0.1")
	// fmt.Println("rootCert\n", string(rootCertPEM))

	// ip := ut.GetIP()
	// internalReg.RegisterInter(tls.LoadX509KeyPair("rootCerts/rootCa.crt", "rootCerts/rootCa.key"), , "India", "sammoRide-First", "West Bengal", ip)
	// network.StartOrederServer("127.0.0.1:2000", "rootCerts/rootCa.crt", "interCerts/sammoride/orderer/1/com/interCa.crt", "interCerts/sammoride/orderer/1/com/interCa.key")
	// register.RegisterPeer("china", "DasRam", "xijing", "127.0.0.1", "interCerts/sammoride/orderer/1/com/interCa.crt", "interCerts/sammoride/orderer/1/com/interCa.key")
	// =======================Cilent
	// fmt.Print(*mode)
	if *mode == "s" {
		network.StartEnrollServer("interCerts/sammoride/orderer/1/com")
	} else {
		client.SendEnrollRequest("India", "Tapas.Das", "west Bengal")
	}

}
