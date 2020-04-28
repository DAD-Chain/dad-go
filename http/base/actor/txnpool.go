package actor

import (
	"time"
	"github.com/dad-go/core/types"
	"github.com/dad-go/eventbus/actor"
	. "github.com/dad-go/errors"
	. "github.com/dad-go/common"
	tc "github.com/dad-go/txnpool/common"
	"errors"
	"fmt"
)

var txnPid *actor.PID
var txnPoolPid *actor.PID

func SetTxActor(actr *actor.PID) {
	txnPid = actr
}
func SetTxnPoolActor(actr *actor.PID) {
	txnPoolPid = actr
}
func AppendTxToPool(txn *types.Transaction) ErrCode {
	future := txnPid.RequestFuture(txn, 10*time.Second)
	result, err := future.Result()
	if err != nil {
		return ErrUnknown
	}
	if errCode, ok := result.(ErrCode); !ok {
		return errCode
	} else {
		return ErrUnknown
	}
}

func GetTxsFromPool(byCount bool) (map[Uint256]*types.Transaction, Fixed64) {
	future := txnPoolPid.RequestFuture(&tc.GetTxnPoolReq{byCount}, 10*time.Second)
	result, err := future.Result()
	if err != nil {
		return nil, 0
	}
	txpool, ok := result.(*tc.GetTxnPoolRsp);
	if !ok {
		return nil, 0
	}
	txMap := make(map[Uint256]*types.Transaction)
	var networkFeeSum Fixed64
	for _, v := range txpool.TxnPool {
		txMap[v.Tx.Hash()] = v.Tx
		networkFeeSum += v.Fee
	}
	return txMap, networkFeeSum

}

func GetTxFromPool(hash Uint256) (tc.TXEntry, error) {

	future := txnPid.RequestFuture(&tc.GetTxnReq{hash}, 10*time.Second)
	result, err := future.Result()
	if err != nil {
		return tc.TXEntry{}, err
	}
	txn, ok := result.(*tc.GetTxnRsp);
	if !ok {
		return tc.TXEntry{}, errors.New("fail")
	}
	if txn == nil {
		return tc.TXEntry{}, errors.New("fail")
	}

	future = txnPid.RequestFuture(&tc.GetTxnStatusReq{hash}, 10*time.Second)
	result, err = future.Result()
	if err != nil {
		return tc.TXEntry{}, err
	}
	txStatus, ok := result.(*tc.GetTxnStatusRsp);
	if !ok {
		return tc.TXEntry{}, errors.New("fail")
	}

	txnEntry := tc.TXEntry{txn.Txn,0,txStatus.TxStatus}

	return txnEntry, nil
}

func GetTxnCnt() ([]uint64, error) {
	future := txnPid.RequestFuture(&tc.GetTxnStats{}, 10*time.Second)
	result, err := future.Result()
	if err != nil {
		return []uint64{}, err
	}
	txnCnt, ok := result.(*tc.GetTxnStatsRsp);
	if !ok {
		return []uint64{}, errors.New("fail")
	}
	fmt.Println(*txnCnt.Count)
	return *txnCnt.Count, nil
}
