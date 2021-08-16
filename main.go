package main

import (
	"flag"
	"fmt"
	"sammoRide/network"
)

func main() {
	fmt.Println("generationg CAs")
	mode := flag.String("mode", "Mode", "Mode of the application")
	flag.Parse()

	fmt.Print(*mode)
	// cert, err := ioutil.ReadFile("rootCerts/rootCa.crt")
	// key, err := ioutil.ReadFile("rootCerts/rootCa.key")
	// ut.CheckErr(err, "main")

	// internalReg.RegisterInter(ut.LoadCertificate(cert), ut.LoadPrivateKey(key), "India", "orderer", "west", "127.0.0.1", "kolkata", "100025")

	// go network.StartEnrollServer("interCerts/orderer")
	// network.StartOrederServer("127.0.0.1", "rootCerts/rootCa.crt",
	// 	"interCerts/orderer/interCa.crt",
	// 	"interCerts/orderer/interCa.key")

	// rootca, rootP, _ := internalReg.RegisterRoot("Presidential")
	// secCa, secP, _ := internalReg.RegisterInter(rootca, rootP, "India", "sammoride", "", "", "", "", "Secretary")
	// ut.VerifyOrderer(rootca, secCa)
	// direcCa, direcP, _ := internalReg.RegisterInter(secCa, secP, "India", "sammoride", "westbengal", "", "kolkata", "700028", "Director")
	// ut.VerifyOrderer(secCa, direcCa)
	// ut.VerifyPeer(rootca, secCa, direcCa)
	// ordererca, _, _ := internalReg.RegisterInter(direcCa, direcP, "India", "sammoride", "west bengal", ut.GetIP(), "kolkata", "700028", "Orderer")
	// ut.VerifyOrderer(direcCa, ordererca)
	// ut.VerifyPeer(secCa, direcCa, ordererca)
	// network.InitFileStructure("OrdererInfo/orderer.yml")
	network.StartService("OrdererInfo/orderer.yml", true)

}
