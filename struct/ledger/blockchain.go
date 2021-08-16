package ledger

import (
	"sammoRide/struct/common"
	"sammoRide/struct/orderer"
	"sammoRide/ut"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type Blockchain struct {
	LastHeader common.BlockHeader
	Database   *badger.DB
}

func (chain *Blockchain) InitBlockchain(filepath string) {
	chaincodeActionPayload := []byte("This is the genesisBlock")
	header := common.Header{ChannelHeader: &common.ChannelHeader{Timestamp: time.Now(), ChannelId: "sammonRide", TxId: "Transaction", Epoch: 0},
		SignatureHeader: &common.SignatureHeader{Driver: []byte("Everyone has equal rights."), Traveler: []byte("Everyone deserves Equal Opportunity")}}

	transactionAction := orderer.TransactionAction{Header: header.Serialize(),
		Payload: chaincodeActionPayload}
	transaction := orderer.Transaction{Actions: []*orderer.TransactionAction{&transactionAction}}

	data := [][]byte{transaction.Serialize()}
	blockdata := common.BlockData{Data: data}

	blockHeader := common.BlockHeader{Number: 1, PreviousHash: nil, DataHash: ut.Hash(blockdata.Serialize())}
	block := common.Block{Header: &blockHeader, Data: &blockdata}

	chain.LastHeader = blockHeader

	db, err := badger.Open(badger.DefaultOptions(filepath))
	ut.CheckErr(err, "NewDatabase/chainDatabase")

	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Set(blockHeader.Serialize(), block.Serialize())
		return err
	})
	ut.CheckErr(err, "InitBlockchain/update")
	chain.Database = db
}

func (chain *Blockchain) Close() {
	chain.Database.Close()
}

func (chain *Blockchain) Update(block common.Block) {
	err := chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.GetHeader().Serialize(), block.Serialize())
		err = txn.Set([]byte("LastHeader"), block.Header.Serialize())
		return err
	})
	ut.CheckErr(err, "ChainDatabase/Update")
	chain.LastHeader = *block.Header
}

func LoadDatabase(filepath string) *Blockchain {
	db, err := badger.Open(badger.DefaultOptions(filepath))
	ut.CheckErr(err, "NewDatabase/chainDatabase")
	chain := Blockchain{}
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("LastHeader"))
		err = item.Value(func(val []byte) error {
			chain.LastHeader = *common.DeSerializeBlockHeader(val)
			return err
		})
		return err
	})
	ut.CheckErr(err, "LoadDatabase")
	chain.Database = db
	return &chain
}
