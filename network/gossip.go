package network

import (
	"bytes"
	"encoding/gob"
	"sammoRide/ut"
)

type GossipData struct {
	Header string
	Data   []byte
}

func (ra *GossipData) gossipSerialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(ra)

	ut.CheckErr(err, "ContactSer/encode")

	return res.Bytes()
}

func GossipDeserialize(data []byte) *GossipData {
	var gData GossipData

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&gData)

	ut.CheckErr(err, "RAD/decode")

	return &gData
}
