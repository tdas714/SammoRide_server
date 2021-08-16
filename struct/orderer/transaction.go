package orderer

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"sammoRide/ut"
)

// The transaction to be sent to the ordering service. A transaction contains
// one or more TransactionAction. Each TransactionAction binds a proposal to
// potentially multiple actions. The transaction is atomic meaning that either
// all actions in the transaction will be committed or none will.  Note that
// while a Transaction might include more than one Header, the Header.creator
// field must be the same in each.
// A single client is free to issue a number of independent Proposal, each with
// their header (Header) and request payload (ChaincodeProposalPayload).  Each
// proposal is independently endorsed generating an action
// (ProposalResponsePayload) with one signature per Endorser. Any number of
// independent proposals (and their action) might be included in a transaction
// to ensure that they are treated atomically.
type Transaction struct {
	// The payload is an array of TransactionAction. An array is necessary to
	// accommodate multiple actions per transaction
	Actions []*TransactionAction

	isValid bool
}

func (m *Transaction) VerifySignatures() bool {
	for _, ta := range m.GetActions() {
		chaincodeAction := DeSerializeChaincodeActionPayload(ta.GetPayload())
		signedProp := DeSerializeSignedProposal(bytes.NewBuffer(chaincodeAction.GetChaincodeProposalPayload()))
		travelerSig := DeSerializeSig(signedProp.TravelerSignature)
		driverSig := DeSerializeSig(signedProp.DriverSignature)

		travelerV := ecdsa.Verify(ut.Keydecode(signedProp.TravelerPublicKey), ut.Hash(signedProp.GetProposalBytes()), travelerSig.R, travelerSig.S)
		driverV := ecdsa.Verify(ut.Keydecode(signedProp.DriverPublicKey), ut.Hash(signedProp.GetProposalBytes()), driverSig.R, driverSig.S)
		if driverV && travelerV {
			m.isValid = true
			return true
		} else {
			m.isValid = false
			return false
		}
	}
	return true
}

func (m *Transaction) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "Transaction/Serialize")
	}
	return js
}

func DeSerializeTransaction(data io.Reader) *Transaction {
	var m *Transaction
	json.NewDecoder(data).Decode(&m)
	return m
}

func (m *Transaction) GetActions() []*TransactionAction {
	if m != nil {
		return m.Actions
	}
	return nil
}

// TransactionAction binds a proposal to its action.  The type field in the
// header dictates the type of action to be applied to the ledger.
type TransactionAction struct {
	// The header of the proposal action, which is the proposal header
	Header []byte
	// The payload of the action as defined by the type in the header For
	// chaincode, it's the bytes of ChaincodeActionPayload
	Payload []byte
}

func (m *TransactionAction) GetHeader() []byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *TransactionAction) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

// ChaincodeActionPayload is the message to be used for the TransactionAction's
// payload when the Header's type is set to CHAINCODE.  It carries the
// chaincodeProposalPayload and an endorsed action to apply to the ledger.
type ChaincodeActionPayload struct {
	// This field contains the bytes of the ChaincodeProposalPayload message from
	// the original invocation (essentially the arguments) after the application
	// of the visibility function. The main visibility modes are "full" (the
	// entire ChaincodeProposalPayload message is included here), "hash" (only
	// the hash of the ChaincodeProposalPayload message is included) or
	// "nothing".  This field will be used to check the consistency of
	// ProposalResponsePayload.proposalHash.  For the CHAINCODE type,
	// ProposalResponsePayload.proposalHash is supposed to be H(ProposalHeader ||
	// f(ChaincodeProposalPayload)) where f is the visibility function.
	ChaincodeProposalPayload []byte
	// The list of actions to apply to the ledger
	Action *ChaincodeEndorsedAction
}

func (m *ChaincodeActionPayload) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "ChainCodActionpayload/Serialize")
	}
	return js
}

func DeSerializeChaincodeActionPayload(data []byte) *ChaincodeActionPayload {
	var m *ChaincodeActionPayload
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *ChaincodeActionPayload) GetChaincodeProposalPayload() []byte {
	if m != nil {
		return m.ChaincodeProposalPayload
	}
	return nil
}

func (m *ChaincodeActionPayload) GetAction() *ChaincodeEndorsedAction {
	if m != nil {
		return m.Action
	}
	return nil
}

// ChaincodeEndorsedAction carries information about the endorsement of a
// specific proposal
type ChaincodeEndorsedAction struct {
	// This is the bytes of the ProposalResponsePayload message signed by the
	// endorsers.  Recall that for the CHAINCODE type, the
	// ProposalResponsePayload's extenstion field carries a ChaincodeAction
	ProposalResponsePayload []byte
	// The endorsement of the proposal, basically the endorser's signature over
	// proposalResponsePayload
	Endorsements []*Endorsement
}

func (m *ChaincodeEndorsedAction) GetProposalResponsePayload() []byte {
	if m != nil {
		return m.ProposalResponsePayload
	}
	return nil
}

func (m *ChaincodeEndorsedAction) GetEndorsements() []*Endorsement {
	if m != nil {
		return m.Endorsements
	}
	return nil
}

// An endorsement is a signature of an endorser over a proposal response.  By
// producing an endorsement message, an endorser implicitly "approves" that
// proposal response and the actions contained therein. When enough
// endorsements have been collected, a transaction can be generated out of a
// set of proposal responses.  Note that this message only contains an identity
// and a signature but no signed payload. This is intentional because
// endorsements are supposed to be collected in a transaction, and they are all
// expected to endorse a single proposal response/action (many endorsements
// over a single proposal response)
type Endorsement struct {
	// Identity of the endorser (e.g. its certificate)
	Endorser []byte
	// Signature of the payload included in ProposalResponse concatenated with
	// the endorser's certificate; ie, sign(ProposalResponse.payload + endorser)
	Signature *Sig
	PublicKey string
}

func (m *Endorsement) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "Endorsement/Serialize")
	}
	return js
}

func DeSerializeEndorsement(data io.Reader) *Endorsement {
	var m *Endorsement
	json.NewDecoder(data).Decode(&m)
	return m
}

func (m *Endorsement) GetEndorser() []byte {
	if m != nil {
		return m.Endorser
	}
	return nil
}

func (m *Endorsement) GetSignature() *Sig {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Endorsement) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

type Sig struct {
	R *big.Int
	S *big.Int
}

func (m *Sig) GetR() *big.Int {
	return m.R
}

func (m *Sig) GetS() *big.Int {
	return m.S
}

func (m *Sig) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "sig/Serialize")
	}
	return js
}

func DeSerializeSig(data []byte) *Sig {
	var m *Sig
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

// ProposalResponsePayload is the payload of a proposal response.  This message
// is the "bridge" between the client's request and the endorser's action in
// response to that request. Concretely, for chaincodes, it contains a hashed
// representation of the proposal (proposalHash) and a representation of the
// chaincode state changes and events inside the extension field.
type ProposalResponsePayload struct {
	// Hash of the proposal that triggered this response. The hash is used to
	// link a response with its proposal, both for bookeeping purposes on an
	// asynchronous system and for security reasons (accountability,
	// non-repudiation). The hash usually covers the entire Proposal message
	// (byte-by-byte).
	ProposalHash []byte
	// Extension should be unmarshaled to a type-specific message. The type of
	// the extension in any proposal response depends on the type of the proposal
	// that the client selected when the proposal was initially sent out.  In
	// particular, this information is stored in the type field of a Header.  For
	// chaincode, it's a ChaincodeAction message
	Extension *ChaincodeAction
}

func (m *ProposalResponsePayload) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "ProposalResponsePayload/Serialize")
	}
	return js
}

func DeSerializeProposalResponsePayload(data []byte) *ProposalResponsePayload {
	var m *ProposalResponsePayload
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m
}

func (m *ProposalResponsePayload) GetProposalHash() []byte {
	if m != nil {
		return m.ProposalHash
	}
	return nil
}

func (m *ProposalResponsePayload) GetExtension() *ChaincodeAction {
	if m != nil {
		return m.Extension
	}
	return nil
}

type SignedProposal struct {
	// The bytes of Proposal
	ProposalBytes []byte
	// Signaure over proposalBytes; this signature is to be verified against
	// the creator identity contained in the header of the Proposal message
	// marshaled as proposalBytes
	DriverSignature   []byte
	TravelerSignature []byte
	DriverPublicKey   string
	TravelerPublicKey string
}

func (m *SignedProposal) Serialize() []byte {
	js, err := json.Marshal(m)
	if err != nil {
		log.Panic(err.Error() + " - " + "signedProposal/Serialize")
	}
	return js
}

func DeSerializeSignedProposal(data io.Reader) *SignedProposal {
	var m *SignedProposal
	json.NewDecoder(data).Decode(&m)
	return m
}

func (m *SignedProposal) GetProposalBytes() []byte {
	if m != nil {
		return m.ProposalBytes
	}
	return nil
}

func (m *SignedProposal) GetDriverSignature() *Sig {
	if m != nil {
		return DeSerializeSig(m.DriverSignature)
	}
	return nil
}

func (m *SignedProposal) GetTravelerSignature() *Sig {
	if m != nil {
		return DeSerializeSig(m.DriverSignature)
	}
	return nil
}
