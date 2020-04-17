package ledger

import (
	. "github.com/dad-go/common"
	"github.com/dad-go/core/states"
	"github.com/dad-go/core/types"
	"github.com/dad-go/crypto"
)

// ILedgerStore provides func with store package.
type ILedgerStore interface {
	//TODO: define the state store func
	SaveBlock(b *types.Block, ledger *Ledger) error
	GetBlock(hash Uint256) (*types.Block, error)
	BlockInCache(hash Uint256) bool
	GetBlockHash(height uint32) (Uint256, error)
	InitLedgerStore(ledger *Ledger) error
	IsDoubleSpend(tx *types.Transaction) bool

	AddHeaders(headers []types.Header, ledger *Ledger) error
	GetHeader(hash Uint256) (*types.Header, error)

	GetTransaction(hash Uint256) (*types.Transaction, error)
	GetTransactionWithHeight(hash Uint256) (*types.Transaction, uint32, error)

	GetAsset(hash Uint256) (*states.AssetState, error)
	GetContract(hash Uint160) (*states.ContractState, error)
	GetAccount(programHash Uint160) (*states.AccountState, error)

	GetCurrentBlockHash() Uint256
	GetCurrentHeaderHash() Uint256
	GetHeaderHeight() uint32
	GetHeight() uint32
	GetHeaderHashByHeight(height uint32) Uint256
	GetBlockRootWithNewTxRoot(txRoot Uint256) Uint256

	GetBookKeeperList() ([]*crypto.PubKey, []*crypto.PubKey, error)
	InitLedgerStoreWithGenesisBlock(genesisblock *types.Block, defaultBookKeeper []*crypto.PubKey) (uint32, error)

	GetQuantityIssued(assetid Uint256) (Fixed64, error)

	//GetUnspent(txid Uint256, index uint16) (*utxo.TxOutput, error)
	//ContainsUnspent(txid Uint256, index uint16) (bool, error)
	//GetUnspentFromProgramHash(programHash Uint160, assetid Uint256) ([]*utxo.UTXOUnspent, error)
	//GetAssets() map[Uint256]*states.AssetState

	IsTxHashDuplicate(txhash Uint256) bool
	IsBlockInStore(hash Uint256) bool
	Close()

	GetCurrentStateRoot() Uint256
	GetIdentity(ontId []byte) ([]byte, error)

	GetStorageItem(key *states.StorageKey) (*states.StorageItem, error)

	GetSysFeeAmount(hash Uint256) (Fixed64, error)
	//GetVotesAndEnrollments(txs []*types.Transaction) ([]*states.VoteState, []*crypto.PubKey, error)
}
