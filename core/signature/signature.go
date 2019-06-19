package signature

 import (
	"dad-go/common"
)

 //SignableData describe the data need be signed.
type SignableData interface {

 	//Get the the SignableData's program hashes
	GetProgramHashes() ([]common.Uint160, error)
}