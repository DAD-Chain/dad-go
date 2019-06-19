package transaction

 import (
	"dad-go/common"
	"dad-go/core/contract"
)


 type BalanceTxInput struct {
	AssetID common.Uint256
	Value common.Fixed8
	Address contract.Address
} 
