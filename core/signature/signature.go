package signature

import (
	"bytes"
	"crypto/sha256"
	"io"

	"github.com/dad-go/common"
	"github.com/dad-go/core/contract/program"
	"github.com/dad-go/crypto"
	. "github.com/dad-go/errors"
	"github.com/dad-go/vm/neovm/interfaces"
)

//SignableData describe the data need be signed.
type SignableData interface {
	interfaces.ICodeContainer

	////Get the the SignableData's program hashes
	GetProgramHashes() ([]common.Address, error)

	SetPrograms([]*program.Program)

	GetPrograms() []*program.Program

	//TODO: add SerializeUnsigned
	SerializeUnsigned(io.Writer) error
}

func SignBySigner(data SignableData, signer Signer) ([]byte, error) {
	return sign(data, signer.PrivKey())
}

func getHashData(data SignableData) []byte {
	buf := new(bytes.Buffer)
	data.SerializeUnsigned(buf)
	return buf.Bytes()
}

func sign(data SignableData, privKey []byte) ([]byte, error) {
	temp := sha256.Sum256(getHashData(data))
	hash := sha256.Sum256(temp[:])

	signature, err := crypto.Sign(privKey, hash[:])
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Signature],Sign failed.")
	}
	return signature, nil
}
