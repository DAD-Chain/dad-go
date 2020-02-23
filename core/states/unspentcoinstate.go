package states

import (
	"io"
	"github.com/dad-go/common/serialization"
)

type UnspentCoinState struct {
	StateBase
	Item []CoinState
}

func (this *UnspentCoinState) Serialize(w io.Writer) error {
	this.StateBase.Serialize(w)
	serialization.WriteUint32(w, uint32(len(this.Item)))
	for _, v := range this.Item {
		serialization.WriteByte(w, byte(v))
	}
	return nil
}

func (this *UnspentCoinState) Deserialize(r io.Reader) error {
	if this == nil {
		this = new(UnspentCoinState)
	}
	err := this.StateBase.Deserialize(r)
	if err != nil {
		return err
	}
	n, err := serialization.ReadUint32(r)
	if err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		this.Item = append(this.Item, CoinState(i))
	}
	return nil
}
