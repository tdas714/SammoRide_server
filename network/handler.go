package network

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"sammoRide/ca"
	"sammoRide/database"
	"sammoRide/struct/common"
	"sammoRide/struct/orderer"
	"sammoRide/ut"
	"sammoRide/ut/client"

	"github.com/dgraph-io/badger/v3"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	ut.CheckErr(err, "HeloRequest")
	fmt.Println("Hello from ", string(bodyBytes))
	// Write "Hello, world!" to the response body
	io.WriteString(w, "Hello, world!\n")
}

func OrdererEnrollHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
	enrollReq := orderer.DeSerializeEnrollDataRequest(r.Body)

	serialNum, err := ioutil.ReadFile(strings.Split(db.UtilsPath, "/")[0] + ut.SERIAL_LOG)
	ut.CheckErr(err, "OrdererEnrollHandler")

	fmt.Println("Enroll Request from ", enrollReq.IpAddr, enrollReq.ListingPort)

	cert, err := ioutil.ReadFile(db.Certificatepath)
	rcert, err := ioutil.ReadFile(db.InterCaPath)
	priv, err := ioutil.ReadFile(db.KeyPath)

	sha := sha1.New()
	sha.Write(enrollReq.Serialize())

	pCert, pPriv := ca.GenDCA("Orderer", ut.LoadCertificate(cert),
		ut.LoadPrivateKey(priv),
		enrollReq.Country, enrollReq.Name, enrollReq.IpAddr, enrollReq.Province, enrollReq.City, enrollReq.PostalCode,
		int64(binary.BigEndian.Uint64(serialNum)+1), sha.Sum(nil),
	)

	pPrivByte, err := x509.MarshalECPrivateKey(pPriv)
	ut.CheckErr(err, "handleReq/pPriveByte")

	enrollReq.PrivateKey = pPrivByte

	db.InsertOrderer(ut.GetBytes(enrollReq.Name+":"+enrollReq.Country+":"+enrollReq.Province+":"+enrollReq.City),
		ut.GetBytes(enrollReq)) //Later will be e-mail

	pCertPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: pCert,
	})

	enrollRes := orderer.EnrollDataResponse{Header: ut.ENROLL_RES,
		IpAddr:      enrollReq.IpAddr,
		PeerCert:    pCertPem,
		PrivateKey:  pPrivByte,
		SenderCert:  cert,
		RootCert:    rcert,
		OrdererList: db.GetRandomOrderer(1),
		PeerList:    []string{}}

	w.Header().Set("Content-Type", "application/json")
	w.Write(enrollRes.Serialize(w))

	err = ioutil.WriteFile(db.UtilsPath+ut.SERIAL_LOG,
		ut.IntToByteArray(int64(binary.BigEndian.Uint64(serialNum)+1)), 0700)
	ut.CheckErr(err, "WriteSerialNumber")
}

func RiderAHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
	// bodyBytes, err := ioutil.ReadAll(r.Body)
	// ut.CheckErr(err, "RiderAHandler")

	riderA := common.RADeserialize(r.Body) //We will need this
	db.UpdateInterestedRider(riderA.Info.IP+":"+riderA.Info.Port, riderA)
	db.Gossip(riderA, 2, "Announcement/rider",
		riderA.Header, riderA.Info.IP, riderA.Info.Port)

	fmt.Print("Rider Announcment from: ", riderA.Info.IP+":"+riderA.Info.Port) //This will change
}

func TransactionCommitmentHandle(w http.ResponseWriter, r *http.Request, db *database.Database) {
	var timeStamp time.Time
	transaction := orderer.DeSerializeTransaction(r.Body)
	for _, ta := range transaction.GetActions() {
		header := common.DeSerializeHeader(ta.GetHeader())
		timeStamp = header.GetChannelHeader().GetTimestamp()

		db.PendingTxs[timeStamp.Unix()] = transaction
	}
	db.PendingTxsKey = append(db.PendingTxsKey, timeStamp.Unix())
	sort.Slice(db.PendingTxsKey, func(i, j int) bool { return db.PendingTxsKey[i] < db.PendingTxsKey[j] })

	if len(db.PendingTxs) >= 1 {
		var data [][]byte
		var valid bool
		keys := db.PendingTxsKey[:1]
		db.PendingTxsKey = db.PendingTxsKey[1:]
		for _, k := range keys {
			valid = db.PendingTxs[k].VerifySignatures()
			data = append(data, db.PendingTxs[k].Serialize())
		}
		blockData := common.BlockData{Data: data}

		// Attach blockChain with database, with saving into badger.db ieithin the Struct
		blockHeader := common.BlockHeader{Number: db.BlockChain.LastHeader.Number + 1, PreviousHash: db.BlockChain.LastHeader.DataHash, DataHash: ut.Hash(blockData.Serialize())}
		block := common.Block{Header: &blockHeader, Data: &blockData}
		if valid {
			db.WorldState.UpdateBlock(&blockData, int(db.BlockChain.LastHeader.GetNumber()))
		}
		// Have to update blockchain
		// if ok {
		db.BlockChain.Update(block)
		plist := db.GetRandomPeer(int64(1))
		for _, peer := range plist {
			if !strings.Contains(peer, "127.0.0.1"+":"+"8443") {
				fmt.Println("Send Gossip: ", peer)
				client.SendData(peer, db.InterCaPath, db.Certificatepath,
					db.KeyPath, "Committer/Block", 1, block.Serialize())
			}
		}

		// }

	}
}

func SnapshotHandler(w http.ResponseWriter, resp *http.Request, db *database.Database) {
	numSnap := ut.ByteArrayToInt(ut.StreamToByte(resp.Body))
	var outBlocks []*common.Block
	err := db.BlockChain.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = int(numSnap)
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			// k := item.Key()
			err := item.Value(func(v []byte) error {
				outBlocks = append(outBlocks, common.DeSerializeBlock(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	snapBlocks := common.SnapshotBlocks{Blocks: outBlocks}

	keyPem, err := ioutil.ReadFile(db.KeyPath)
	ut.CheckErr(err, "snapBlocks/Keypem")

	r, s, err := ecdsa.Sign(rand.Reader, ut.LoadPrivateKey(keyPem), ut.Hash(snapBlocks.Serialize()))
	ut.CheckErr(err, "snapblocks/sign")
	sig := orderer.Sig{R: r, S: s}

	snapEnv := common.SnapshotEnvelop{Data: snapBlocks.Serialize(), Signature: &sig, PublicKey: ut.Keyencode(&ut.LoadPrivateKey(keyPem).PublicKey)}

	w.Header().Set("Content-Type", "application/json")
	w.Write(snapEnv.Serialize())
}
