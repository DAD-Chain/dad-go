package payload

import (
	"github.com/dad-go/crypto"
	"io"
)


type Enrollment struct {
	PublicKey *crypto.PubKey
}

func (e *Enrollment) Data(version byte) []byte {
	return []byte{0}
}

func (e *Enrollment) Serialize(w io.Writer, version byte) error {
	if err := e.PublicKey.Serialize(w); err != nil {
		return err
	}
	return nil
}

func (e *Enrollment) Deserialize(r io.Reader, version byte) error {
	pk := new(crypto.PubKey)
	if err := pk.DeSerialize(r); err != nil {
		return err
	}
	e.PublicKey = pk
	return nil
}

