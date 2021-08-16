package orderer

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type EnrollDataRequest struct {
	Country     string
	Name        string
	Province    string
	IpAddr      string
	City        string
	PostalCode  string
	ListingPort string
	PrivateKey  []byte
}

func (m *EnrollDataRequest) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "EnrollDataRequest/Serialize")
	}
	return js
}

func DeSerializeEnrollDataRequest(data io.Reader) *EnrollDataRequest {
	var m *EnrollDataRequest
	json.NewDecoder(data).Decode(&m)
	return m
}

type EnrollDataResponse struct {
	Header      string
	IpAddr      string
	PeerCert    []byte
	PrivateKey  []byte
	SenderCert  []byte
	RootCert    []byte
	PeerList    []string
	OrdererList []string
}

func (m *EnrollDataResponse) Serialize(rw http.ResponseWriter) []byte {
	js, err := json.Marshal(m)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return []byte{}
	}
	return js
}

func DeSerializeEnrollDataResponse(data io.Reader) *EnrollDataResponse {
	var m *EnrollDataResponse
	json.NewDecoder(data).Decode(&m)
	return m
}
