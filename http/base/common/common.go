package common

import (
	. "github.com/dad-go/common"
	"github.com/dad-go/common/log"
	"github.com/dad-go/core/types"
	. "github.com/dad-go/errors"
	. "github.com/dad-go/http/base/actor"
	"strconv"
)

//multiplexer that keeps track of every function to be called on specific rpc call

type TxAttributeInfo struct {
	Usage types.TransactionAttributeUsage
	Data  string
}

type AmountMap struct {
	Key   Uint256
	Value Fixed64
}

type ProgramInfo struct {
	Code      string
	Parameter string
}

type Fee struct {
	Amount Fixed64
	Payer  string
}
type PubKeyInfo struct {
	X, Y string
}
type Sig struct {
	PubKeys []PubKeyInfo
	M       uint8
	SigData []string
}
type Transactions struct {
	Version    byte
	Nonce      uint32
	TxType     types.TransactionType
	Payload    PayloadInfo
	Attributes []TxAttributeInfo
	Fee        []Fee
	NetworkFee string
	Sigs       []Sig
	Hash       string
}

type BlockHead struct {
	Version          uint32
	PrevBlockHash    string
	TransactionsRoot string
	BlockRoot        string
	StateRoot        string
	Timestamp        uint32
	Height           uint32
	ConsensusData    uint64
	NextBookKeeper   string

	BookKeepers []string
	SigData     []string

	Hash string
}

type BlockInfo struct {
	Hash         string
	BlockData    *BlockHead
	Transactions []*Transactions
}

type NodeInfo struct {
	NodeState   uint   // node status
	NodePort    uint16 // The nodes's port
	ID          uint64 // The nodes's id
	NodeTime    int64
	NodeVersion uint32   // The network protocol the node used
	NodeType    uint64   // The services the node supplied
	Relay       bool     // The relay capability of the node (merge into capbility flag)
	Height      uint32   // The node latest block height
	TxnCnt      []uint64 // The transactions be transmit by this node
	//RxTxnCnt uint64 // The transaction received by this node
}

type ConsensusInfo struct {
	// TODO
}

type TXNAttrInfo struct {
	Height  uint32
	Type    int
	ErrCode int
}

type TXNEntryInfo struct {
	Txn   Transactions  // transaction which has been verified
	Fee   int64         // Total fee per transaction
	Attrs []TXNAttrInfo // the result from each validator
}

func TransArryByteToHexString(ptx *types.Transaction) *Transactions {
	trans := new(Transactions)
	trans.TxType = ptx.TxType
	trans.Payload = TransPayloadToHex(ptx.Payload)

	trans.Attributes = make([]TxAttributeInfo, len(ptx.Attributes))
	for i, v := range ptx.Attributes {
		trans.Attributes[i].Usage = v.Usage
		trans.Attributes[i].Data = ToHexString(v.Data)
	}
	for _, fee := range ptx.Fee {
		e := Fee{fee.Amount, fee.Payer.ToHexString()}
		trans.Fee = append(trans.Fee, e)
	}
	for _, sig := range ptx.Sigs {
		e := Sig{M: sig.M}
		for i := 0; i < len(sig.PubKeys); i++ {
			e.PubKeys = append(e.PubKeys, PubKeyInfo{sig.PubKeys[i].X.String(), sig.PubKeys[i].Y.String()})
			e.SigData = append(e.SigData, ToHexString(sig.SigData[i]))
		}
		trans.Sigs = append(trans.Sigs, e)
	}
	networkfee := ptx.GetNetworkFee()
	trans.NetworkFee = strconv.FormatInt(int64(networkfee), 10)

	mhash := ptx.Hash()
	trans.Hash = ToHexString(mhash.ToArray())
	return trans
}

func VerifyAndSendTx(txn *types.Transaction) ErrCode {
	// if transaction is verified unsucessfully then will not put it into transaction pool
	if errCode := AppendTxToPool(txn); errCode != ErrNoError {
		log.Warn("Can NOT add the transaction to TxnPool")
		log.Info("[httpjsonrpc] VerifyTransaction failed when AppendTxnPool.")
		return errCode
	}

	if err := Xmit(txn); err != nil {
		log.Error("Xmit Tx Error:Xmit transaction failed.", err)
		return ErrXmitFail
	}
	return ErrNoError
}
