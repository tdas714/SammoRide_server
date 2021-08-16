package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sammoRide/ut"

	"gopkg.in/yaml.v3"
)

type InputInfo struct {
	OrdererInfo *OrdererInfo `yaml:"Info"`
	Path        *FilePath    `yaml:"Path"`
}

type FilePath struct {
	CertificatePath string `yaml:"CertificatePath"`
	KeyPath         string `yaml:"KeyPath"`
	CAsPath         string `yaml:"CAsPath"`
	PeerDBPath      string `yaml:"PeerDBPath"`
	OrdererfileDB   string `yaml:"OrdererFileDB"`
	ChainDB         string `yaml:"ChainDB"`
	StatePath       string `yaml:"StatePath"`
}

type OrdererInfo struct {
	Country    string `yaml:"Country"`
	Name       string `yaml:"Name"`
	Province   string `yaml:"Province"`
	City       string `yaml:"City"`
	IP         string `yaml:"IP"`
	Postalcode string `yaml:"PostalCode"`
	Port       string `yaml:"Port"`
	PublicKey  string
}

func (c *InputInfo) Parse(filename string) (*OrdererInfo, *FilePath) {

	yamlFile, err := ioutil.ReadFile(filename)
	ut.CheckErr(err, "YamlFile Get")
	err = yaml.Unmarshal(yamlFile, &c)
	ut.CheckErr(err, " Unmarshal Error")
	return c.OrdererInfo, c.Path

}

type ClientInfo struct {
	Country    string
	Name       string
	Province   string
	City       string
	IP         string
	Postalcode string
	Port       string
	PublicKey  string
}

type RiderAnnouncement struct {
	Header      int64
	Latitude    string
	Longitude   string
	Avalability string
	Info        *ClientInfo
}

func (ra *RiderAnnouncement) RASerialize() []byte {
	js, err := json.Marshal(ra)
	ut.CheckErr(err, "RAS/encode")

	return js
}

func RADeserialize(data io.Reader) *RiderAnnouncement {
	var riderA *RiderAnnouncement

	json.NewDecoder(data).Decode(&riderA)

	return riderA
}
