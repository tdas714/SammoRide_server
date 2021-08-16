package database

import (
	"bytes"
	"fmt"
	"math/rand"

	"sammoRide/struct/common"
	"sammoRide/struct/ledger"
	"sammoRide/struct/orderer"
	"sammoRide/ut"
	"sammoRide/ut/client"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type Database struct {
	Info             *common.OrdererInfo
	PeerDB           *badger.DB
	PeerList         []string
	OrdererDB        *badger.DB
	OrdererList      []string
	InterestedRiders map[string]*common.RiderAnnouncement
	InterCaPath      string
	RootCaPath       string
	Certificatepath  string
	KeyPath          string
	UtilsPath        string
	GossipSentList   map[int64]string
	PendingTxsKey    []int64
	PendingTxs       map[int64]*orderer.Transaction
	BlockChain       *ledger.Blockchain
	WorldState       *ledger.WorldState
}

func (database *Database) Close() {
	database.PeerDB.Close()
	database.OrdererDB.Close()
	database.BlockChain.Close()
	database.WorldState.Close(database.UtilsPath + "/WorldState.json")
}

func NewDatabase(ord *common.OrdererInfo, filepath, ordererPath, intercaPath, rootCaPath, crtPath, keyPath, chainPath, utilsPath string, isGenesis bool) *Database {
	peerdb, err := badger.Open(badger.DefaultOptions(filepath))
	ut.CheckErr(err, "NewDatabase/db")

	var client *orderer.EnrollDataRequest

	var PList, OList []string

	err = peerdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			// k := item.Key()
			err := item.Value(func(v []byte) error {
				client = orderer.DeSerializeEnrollDataRequest(bytes.NewBuffer(v))
				PList = append(PList, client.IpAddr+":"+client.ListingPort)
				return nil
			})
			ut.CheckErr(err, "NewDatabase")
		}
		return nil
	})

	ut.CheckErr(err, "PeerList/Db")

	ordererdb, err := badger.Open(badger.DefaultOptions(ordererPath))
	ut.CheckErr(err, "NewDatabase/db")

	err = ordererdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			// k := item.Key()
			err := item.Value(func(v []byte) error {
				client = orderer.DeSerializeEnrollDataRequest(bytes.NewBuffer(v))
				OList = append(OList, client.IpAddr+":"+client.ListingPort)
				return nil
			})
			ut.CheckErr(err, "NewDatabase")
		}
		return nil
	})

	ut.CheckErr(err, "PeerList/Db")

	iRider := make(map[string]*common.RiderAnnouncement)
	var blockchain ledger.Blockchain

	if isGenesis {
		blockchain = ledger.Blockchain{}
		blockchain.InitBlockchain(chainPath)
	} else {
		blockchain = *ledger.LoadDatabase(chainPath)
	}

	data := Database{Info: ord, PeerDB: peerdb, OrdererDB: ordererdb, PeerList: PList, InterCaPath: intercaPath, RootCaPath: rootCaPath,
		Certificatepath: crtPath, KeyPath: keyPath, InterestedRiders: iRider,
		GossipSentList: make(map[int64]string), PendingTxs: make(map[int64]*orderer.Transaction), BlockChain: &blockchain,
		WorldState: ledger.Init(), UtilsPath: utilsPath}

	return &data
}

func (database *Database) InsertNode(key, value []byte) {
	client := orderer.DeSerializeEnrollDataRequest(bytes.NewBuffer(value))
	err := database.PeerDB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	ut.CheckErr(err, "Database/Update")
	if !Contains(database.PeerList, client.IpAddr+":"+client.ListingPort) {
		database.PeerList = append(database.PeerList,
			client.IpAddr+":"+client.ListingPort)
	}
	fmt.Println("Database nPeer list: ", database.PeerList)
}

func (database *Database) InsertOrderer(key, value []byte) {
	client := orderer.DeSerializeEnrollDataRequest(bytes.NewBuffer(value))
	err := database.OrdererDB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	ut.CheckErr(err, "Database/Update")
	if !Contains(database.OrdererList, client.IpAddr+":"+client.ListingPort) {
		database.OrdererList = append(database.OrdererList,
			client.IpAddr+":"+client.ListingPort)
	}
	fmt.Println("Database nPeer list: ", database.PeerList)
}

func (database *Database) GetRandomPeer(num int64) []string {
	var gList []string
	var selectedPeer string
	i := 1
	for {
		rand.Seed(int64(i) + time.Now().Unix())
		selectedPeer = database.PeerList[rand.Intn(len(database.PeerList))]
		if !Contains(gList, selectedPeer) {
			gList = append(gList, selectedPeer)
		}
		if len(gList) >= int(num) || i >= len(database.PeerList) {
			break
		}
		i++
	}
	return gList
}

func (database *Database) GetRandomOrderer(num int64) []string {
	var gList []string
	var selectedOrderer string
	i := 1
	for {
		if len(gList) >= int(num) || i >= len(database.OrdererList) {
			gList = append(gList, ut.GetIP())
			break
		}
		rand.Seed(int64(i) + time.Now().Unix())
		selectedOrderer = database.OrdererList[rand.Intn(len(database.OrdererList))]
		if !Contains(gList, selectedOrderer) {
			gList = append(gList, selectedOrderer)
		}

		i++
	}
	return gList
}

func (database *Database) Gossip(data *common.RiderAnnouncement, num int, domain string, header int64, rejectedIp, rejectedPort string) {
	_, ok := database.GossipSentList[header]
	if !ok {

		plist := database.GetRandomPeer(int64(num))
		for _, peer := range plist {
			if !strings.Contains(peer, rejectedPort+":"+rejectedPort) {
				fmt.Println("Send Gossip: ", peer)
				client.SendData(peer, database.RootCaPath, database.Certificatepath,
					database.KeyPath, domain, 1, data.RASerialize())
			}
		}
		for _, o := range database.OrdererList {
			if !strings.Contains(o, database.Info.IP+":"+"8443") {

				client.SendData(o, database.RootCaPath, database.Certificatepath,
					database.KeyPath, domain, 1, data.RASerialize())
			}
		}
		database.GossipSentList[header] = data.Info.IP + ":" + data.Info.Port

	}
}

func (database *Database) UpdateInterestedRider(rider string, AnnoucR *common.RiderAnnouncement) {

	_, ok := database.InterestedRiders[rider]
	if ok {
		delete(database.InterestedRiders, rider)
	} else {
		database.InterestedRiders[rider] = AnnoucR
	}
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// =-===========================================================
