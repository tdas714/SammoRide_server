package main

import (
	"fmt"
	// "sammoRide/ca"
	"sammoRide/register/internalReg"
)

func main() {
	fmt.Println("This is a test")
	// _, rootCertPEM, _ := ca.GenCARoot("India", "SammoRide", "127.0.0.1")
	// fmt.Println("rootCert\n", string(rootCertPEM))
	internalReg.RegisterRoot()
}
