package ledger

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"sammoRide/struct/common"
	"sammoRide/struct/orderer"
	"sammoRide/ut"
)

type StateValue struct {
	Value map[string]interface{}
}

func (m *StateValue) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "StateValue/Serialize")
	}
	return js
}

func DeSerializeStateValue(data []byte) *StateValue {
	var m *StateValue
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

type State struct {
	StateValue []byte
	Version    int
}

func (m *State) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "State/Serialize")
	}
	return js
}

func DeSerializeState(data []byte) *State {
	var m *State
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

type WorldState struct {
	CurrentState map[string][]byte
}

func Init() *WorldState {
	value := make(map[string]interface{})
	statev := StateValue{Value: value}
	state := State{StateValue: statev.Serialize(), Version: 0}
	currentState := make(map[string][]byte)
	currentState["Init"] = state.Serialize()
	ws := &WorldState{CurrentState: currentState}
	return ws
}
func (ws *WorldState) Close(filename string) {
	jsonData, err := json.Marshal(ws.CurrentState)
	ut.CheckErr(err, "WorldState/Close")
	jsonFile, err := os.Create(filename)
	ut.CheckErr(err, "WorldState/JsonFile")
	defer jsonFile.Close()
	jsonFile.Write(jsonData)
	jsonFile.Close()

}

func (ws *WorldState) Update(key string, value []byte, isdelete bool) {
	statebytes, ok := ws.CurrentState[key]
	if ok {
		if isdelete {
			ws.Delete(key)
		} else {
			state := DeSerializeState(statebytes)
			state.StateValue = value
			state.Version += 1
			ws.CurrentState[key] = state.Serialize()
		}
	} else {
		state := State{StateValue: value, Version: 0}
		ws.CurrentState[key] = state.Serialize()
	}
}

func (ws *WorldState) Delete(key string) {
	delete(ws.CurrentState, key)
}

func (ws *WorldState) UpdateBlock(blockData *common.BlockData, lastHeight int) bool {

	for _, d := range blockData.Data {
		t := orderer.DeSerializeTransaction(bytes.NewBuffer(d))
		for _, ta := range t.GetActions() {
			chaincodeActionPayload := orderer.DeSerializeChaincodeActionPayload(ta.GetPayload())
			chaincodeAction := orderer.DeSerializeProposalResponsePayload(chaincodeActionPayload.GetAction().GetProposalResponsePayload()).GetExtension()
			kvset := common.DeSerializeKVRWSet(chaincodeAction.GetResults())
			if len(kvset.Reads) != 0 {
				for _, r := range kvset.Reads {
					_, ok := ws.CurrentState[r.GetKey()]
					if !ok || r.GetVersion().BlockNum >= uint64(lastHeight) {
						return false
					}
				}
			}
			if len(kvset.Writes) != 0 {
				for _, w := range kvset.Writes {
					ws.Update(w.GetKey(), w.GetValue(), w.GetIsDelete())
				}
			}

		}
	}
	return true
}
