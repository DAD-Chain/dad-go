package payload

import (
	"github.com/dad-go/common"
	vmtypes "github.com/dad-go/vm/types"
	"io"
	. "github.com/dad-go/errors"
)

type InvokeCode struct {
	GasLimit common.Fixed64
	Code     vmtypes.VmCode
}

func (self *InvokeCode) Serialize(w io.Writer) error {
	var err error
	err = self.GasLimit.Serialize(w)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode GasLimit Serialize failed.")
	}
	err = self.Code.Serialize(w)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode Code Serialize failed.")
	}
	return err
}

func (self *InvokeCode) Deserialize(r io.Reader) error {
	var err error

	err = self.GasLimit.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode GasLimit Deserialize failed.")
	}
	err = self.Code.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode Code Deserialize failed.")
	}
	return nil
}
