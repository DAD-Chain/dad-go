package transaction

 import (
	"dad-go/common"
)


 type BalanceTxInput struct {
	AssetID common.Uint256
	Value common.Fixed8
	ProgramHash common.Uint160
} 
