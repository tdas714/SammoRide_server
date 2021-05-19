package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sammoRide/ut"
)

func SendEnrollRequest(country, name, province string) {
	var enrollRes *ut.PeerEnrollDataResponse

	c, err := net.Dial("tcp", ut.GetIP()+":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	reqEnrroll := ut.PeerEnrollDataRequest{Header: ut.ENROLL_REQ, Country: country, Name: name, Province: province,
		IpAddr: ut.GetIP()}

	b, err := json.Marshal(reqEnrroll)

	fmt.Print(">> ")
	text := string(b)
	fmt.Fprintf(c, text+"\n")

	message, err := bufio.NewReader(c).ReadBytes('\n')
	ut.CheckErr(err)
	json.Unmarshal(message, &enrollRes)

	fmt.Print("->: " + enrollRes.Header)
}
