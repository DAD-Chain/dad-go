package asset

import (
	"github.com/DAD-Chain/dad-go/common"
	"github.com/DAD-Chain/dad-go/common/serialization"
	. "github.com/DAD-Chain/dad-go/errors"
	"errors"
	"io"
)

//AssetType
type AssetType byte

const (
	Currency AssetType = 0x00
	Share    AssetType = 0x01
	Invoice  AssetType = 0x10
	Token    AssetType = 0x11
)

type AssetRecordType byte

//onchain dad-go is planed to support UTXO and Balance
const (
	UTXO    AssetRecordType = 0x00
	Balance AssetRecordType = 0x01
)

//define the asset stucture in onchain dad-go
//registered asset will be assigned to contract address
type Asset struct {
	ID         common.Uint256
	Name       string
	Precision  byte
	AssetType  AssetType
	RecordType AssetRecordType
}

// Serialize is the implement of SignableData interface.
func (a *Asset) Serialize(w io.Writer) error {
	a.ID.Serialize(w)
	err := serialization.WriteVarString(w, a.Name)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], Name serialize failed.")
	}
	_, err = w.Write([]byte{byte(a.Precision)})
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], Precision serialize failed.")
	}
	_, err = w.Write([]byte{byte(a.AssetType)})
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], AssetType serialize failed.")
	}
	_, err = w.Write([]byte{byte(a.RecordType)})
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[Asset], RecordType serialize failed.")
	}
	return nil
}

// Deserialize is the implement of SignableData interface.
func (a *Asset) Deserialize(r io.Reader) error {
	a.ID.Deserialize(r)
	vars, err := serialization.ReadVarString(r)
	if err != nil {
		return NewDetailErr(errors.New("[Asset], Name deserialize failed."), ErrNoCode, "")
	}
	a.Name = vars
	p := make([]byte, 1)
	n, err := r.Read(p)
	if n > 0 {
		a.Precision = p[0]
	} else {
		return NewDetailErr(errors.New("[Asset], Precision deserialize failed."), ErrNoCode, "")
	}
	n, err = r.Read(p)
	if n > 0 {
		a.AssetType = AssetType(p[0])
	} else {
		return NewDetailErr(errors.New("[Asset], AssetType deserialize failed."), ErrNoCode, "")
	}
	n, err = r.Read(p)
	if n > 0 {
		a.RecordType = AssetRecordType(p[0])
	} else {
		return NewDetailErr(errors.New("[Asset], RecordType deserialize failed."), ErrNoCode, "")
	}
	return nil
}

// as the import cycle move to ledger.go
/*func GetAsset(assetId common.Uint256) *Asset {
	fmt.Println("///asset/GetAsset")
	asset, err:= ledger.DefaultLedger.Store.GetAsset(assetId)
	if err != nil {
		return nil,NewDetailErr(err, ErrNoCode, "[Asset], GetAsset failed.")
	}
	return asset, nil
}*/
