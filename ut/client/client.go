package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sammoRide/ut"
)

func SendEnrollRequest(country, name, province string) {
	enrollRes := &ut.PeerEnrollDataRequest{Country: "India", Name: "Tapas", Province: "West Bengal", IpAddr: ut.GetIP()}
	json_data, err := json.Marshal(enrollRes)

	if err != nil {
		log.Fatal(err)
	}
	// var enrollRes *ut.PeerEnrollDataResponse
	path := fmt.Sprintf("http://localhost:8080/post")
	resp, err := http.Post(path, "application/json", bytes.NewBuffer(json_data))
	ut.CheckErr(err)
	println(resp)
}
