package common

import (
	"bytes"
	"encoding/json"
	"log"
)

// KVRWSet encapsulates the read-write set for a chaincode that operates upon a KV or Document data model
// This structure is used for both the public data and the private data
type KVRWSet struct {
	Reads  []*KVRead
	Writes []*KVWrite
}

func (m *KVRWSet) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "KVRWset/Serialize")
	}
	return js
}

func DeSerializeKVRWSet(data []byte) *KVRWSet {
	var m *KVRWSet
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *KVRWSet) GetReads() []*KVRead {
	if m != nil {
		return m.Reads
	}
	return nil
}

func (m *KVRWSet) GetWrites() []*KVWrite {
	if m != nil {
		return m.Writes
	}
	return nil
}

// KVRead captures a read operation performed during transaction simulation
// A 'nil' version indicates a non-existing key read by the transaction
type KVRead struct {
	Key     string
	Version *Version
}

func (m *KVRead) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KVRead) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

// KVWrite captures a write (update/delete) operation performed during transaction simulation
type KVWrite struct {
	Key      string
	IsDelete bool
	Value    []byte
}

func (m *KVWrite) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KVWrite) GetIsDelete() bool {
	if m != nil {
		return m.IsDelete
	}
	return false
}

func (m *KVWrite) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// Version encapsulates the version of a Key
// A version of a committed key is maintained as the height of the transaction that committed the key.
// The height is represenetd as a tuple <blockNum, txNum> where the txNum is the position of the transaction
// (starting with 0) within block
type Version struct {
	BlockNum uint64
	TxId     string
}

func (m *Version) GetBlockNum() uint64 {
	if m != nil {
		return m.BlockNum
	}
	return 0
}

func (m *Version) GetTxNum() string {
	if m != nil {
		return m.TxId
	}
	return "0"
}
