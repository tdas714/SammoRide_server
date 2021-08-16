package common

import (
	"bytes"
	"encoding/json"
	"log"
	"sammoRide/struct/orderer"
	"time"
)

type Header struct {
	ChannelHeader   *ChannelHeader
	SignatureHeader *SignatureHeader
}

func (m *Header) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "Common/Serialize")
	}
	return js
}

func DeSerializeHeader(data []byte) *Header {
	var m *Header
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *Header) GetChannelHeader() *ChannelHeader {
	if m != nil {
		return m.ChannelHeader
	}
	return nil
}

func (m *Header) GetSignatureHeader() *SignatureHeader {
	if m != nil {
		return m.SignatureHeader
	}
	return nil
}

// Header is a generic replay prevention and identity message to include in a signed payload
type ChannelHeader struct {
	// Timestamp is the local time when the message was created
	// by the sender
	Timestamp time.Time

	ChannelId string
	// An unique identifier that is used end-to-end.
	//  -  set by higher layers such as end user or SDK
	//  -  passed to the endorser (which will check for uniqueness)
	//  -  as the header is passed along unchanged, it will be
	//     be retrieved by the committer (uniqueness check here as well)
	//  -  to be stored in the ledger
	TxId string
	// The epoch in which this header was generated, where epoch is defined based on block height
	// Epoch in which the response has been generated. This field identifies a
	// logical window of time. A proposal response is accepted by a peer only if
	// two conditions hold:
	// 1. the epoch specified in the message is the current epoch
	// 2. this message has been only seen once during this epoch (i.e. it hasn't
	//    been replayed)
	Epoch uint64
}

func (m *ChannelHeader) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "Common/Serialize")
	}
	return js
}

func DeSerializeChannelHeader(data []byte) *ChannelHeader {
	var m *ChannelHeader
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *ChannelHeader) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Now().AddDate(25, 0, 0)
}

func (m *ChannelHeader) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *ChannelHeader) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *ChannelHeader) GetEpoch() uint64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

type SignatureHeader struct {
	// Creator of the message, a marshaled msp.SerializedIdentity
	Driver   []byte
	Traveler []byte
	// Arbitrary number that may only be used once. Can be used to detect replay attacks.
	DriverNonce   []byte
	TravelerNonce []byte
}

func (m *SignatureHeader) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "Common/Serialize")
	}
	return js
}

func DeSerializeSignatureHeader(data []byte) *SignatureHeader {
	var m *SignatureHeader
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

// This is finalized block structure to be shared among the orderer and peer
// Note that the BlockHeader chains to the previous BlockHeader, and the BlockData hash is embedded
// in the BlockHeader.  This makes it natural and obvious that the Data is included in the hash, but
// the Metadata is not.
type Block struct {
	Header   *BlockHeader
	Data     *BlockData
	Metadata *BlockMetadata
}

func (m *Block) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "Block/Serialize")
	}
	return js
}

func DeSerializeBlock(data []byte) *Block {
	var m *Block
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *Block) GetHeader() *BlockHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Block) GetData() *BlockData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Block) GetMetadata() *BlockMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

// BlockHeader is the element of the block which forms the block chain
// The block header is hashed using the configured chain hashing algorithm
// over the ASN.1 encoding of the BlockHeader
type BlockHeader struct {
	Number       uint64
	PreviousHash []byte
	DataHash     []byte
}

func (m *BlockHeader) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "BlockHeader/Serialize")
	}
	return js
}

func DeSerializeBlockHeader(data []byte) *BlockHeader {
	var m *BlockHeader
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *BlockHeader) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *BlockHeader) GetPreviousHash() []byte {
	if m != nil {
		return m.PreviousHash
	}
	return nil
}

func (m *BlockHeader) GetDataHash() []byte {
	if m != nil {
		return m.DataHash
	}
	return nil
}

// This is a array of tansactions
type BlockData struct {
	Data [][]byte
}

func (m *BlockData) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "BlockData/Serialize")
	}
	return js
}

func DeSerializeBlockData(data []byte) *BlockData {
	var m *BlockData
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *BlockData) GetData() [][]byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type BlockMetadata struct {
	Metadata [][]byte
}

func (m *BlockMetadata) GetMetadata() [][]byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type SnapshotBlocks struct {
	Blocks []*Block // Chaneg it to [][]byte
}

func (m *SnapshotBlocks) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "SnapshotBlocks/Serialize")
	}
	return js
}

func DeSerializeSnapshotBlocks(data []byte) *SnapshotBlocks {
	var m *SnapshotBlocks
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

type SnapshotEnvelop struct {
	Data      []byte
	Signature *orderer.Sig
	PublicKey string
}

func (m *SnapshotEnvelop) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "SnapshotEnvelop/Serialize")
	}
	return js
}

func DeSerializeSnapshotEnvelop(data []byte) *SnapshotEnvelop {
	var m *SnapshotEnvelop
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}
