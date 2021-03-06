package orderer

import (
	"encoding/json"
	"io"
	"log"
)

// ChaincodeAction contains the executed chaincode results, response, and event.
type ChaincodeAction struct {
	// This field contains the read set and the write set produced by the
	// chaincode executing this invocation.
	Results []byte
	// This field contains the event generated by the chaincode.
	// Only a single marshaled ChaincodeEvent is included.
	Events []byte
	// This field contains the ChaincodeID of executing this invocation. Endorser
	// will set it with the ChaincodeID called by endorser while simulating proposal.
	// Committer will validate the version matching with latest chaincode version.
	// Adding ChaincodeID to keep version opens up the possibility of multiple
	// ChaincodeAction per transaction.
	ChaincodeId *ChaincodeID
}

func (m *ChaincodeAction) GetResults() []byte {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *ChaincodeAction) GetEvents() []byte {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *ChaincodeAction) GetChaincodeId() *ChaincodeID {
	if m != nil {
		return m.ChaincodeId
	}
	return nil
}

//ChaincodeEvent is used for events and registrations that are specific to chaincode
//string type - "chaincode"
type ChaincodeEvent struct {
	ChaincodeId *ChaincodeID
	TxId        string
	EventName   string
	Payload     []byte
}

func (m *ChaincodeEvent) GetChaincodeId() *ChaincodeID {
	if m != nil {
		return m.ChaincodeId
	}
	return nil
}

func (m *ChaincodeEvent) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func (m *ChaincodeEvent) GetEventName() string {
	if m != nil {
		return m.EventName
	}
	return ""
}

func (m *ChaincodeEvent) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type EventsStruct struct {
	ChaincodeEvents []*ChaincodeEvent
}

func (m *EventsStruct) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "Events/Serialize")
	}
	return js
}

func DeSerializeEvetns(data io.Reader) *EventsStruct {
	var m *EventsStruct
	json.NewDecoder(data).Decode(&m)
	return m
}

//ChaincodeID contains the path as specified by the deploy transaction
//that created it as well as the hashCode that is generated by the
//system for the path. From the user level (ie, CLI, REST API and so on)
//deploy transaction is expected to provide the path and other requests
//are expected to provide the hashCode. The other value will be ignored.
//Internally, the structure could contain both values. For instance, the
//hashCode will be set when first generated using the path
type ChaincodeID struct {
	//all other requests will use the name (really a hashcode) generated by
	//the deploy transaction
	Name string
	//user friendly version name for the chaincode
	Version string
}

func (m *ChaincodeID) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ChaincodeID) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}
