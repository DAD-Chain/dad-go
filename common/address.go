package common

import (
	"crypto/sha256"
	"errors"
	"github.com/dad-go/common/log"
	. "github.com/dad-go/errors"
	"io"
	"math/big"

	"github.com/itchyny/base58-go"
	"fmt"
)

const AddrLen int = 20

type Address [AddrLen]byte

func (u *Address) ToArray() []byte {
	x := append([]byte{}, u[:]...)
	return x
}


func (self *Address) ToHexString() string {
	return fmt.Sprintf("%x", self[:])
}

func (self *Address) Serialize(w io.Writer) error {
	_, err := w.Write(self[:])
	return err
}

func (self *Address) Deserialize(r io.Reader) error {
	n, err := r.Read(self[:])
	if n != len(self[:]) || err != nil {
		return errors.New("deserialize Address error")
	}
	return nil
}


func (f *Address) ToBase58() string {
	data := append([]byte{0x41}, f[:]...)
	temp := sha256.Sum256(data)
	temps := sha256.Sum256(temp[:])
	data = append(data, temps[0:4]...)

	bi := new(big.Int).SetBytes(data).String()
	encoded, _ := base58.BitcoinEncoding.Encode([]byte(bi))
	return string(encoded)
}

func Uint160ParseFromBytes(f []byte) (Address, error) {
	if len(f) != AddrLen {
		return Address{}, NewDetailErr(errors.New("[Common]: Uint160ParseFromBytes err, len != 20"), ErrNoCode, "")
	}

	var hash [20]uint8
	for i := 0; i < 20; i++ {
		hash[i] = f[i]
	}
	return Address(hash), nil
}

func AddressFromBase58(encoded string) (Address, error) {
	decoded, err := base58.BitcoinEncoding.Decode([]byte(encoded))
	if err != nil {
		return Address{}, err
	}

	x, _ := new(big.Int).SetString(string(decoded), 10)
	log.Tracef("[ToAddress] x: ", x.Bytes())

	ph, err := Uint160ParseFromBytes(x.Bytes()[1:21])
	if err != nil {
		return Address{}, err
	}

	log.Tracef("[AddressToProgramHash] programhash: %x", ph[:])

	addr := ph.ToBase58()

	log.Tracef("[AddressToProgramHash] encoded: %s", addr)

	if addr != encoded {
		return Address{}, errors.New("[AddressFromBase58]: decode encoded verify failed.")
	}

	return ph, nil
}
