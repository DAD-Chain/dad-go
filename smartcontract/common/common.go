package common

import "github.com/dad-go/vm/neovm/types"

func ConvertTypes(item types.StackItemInterface) interface{} {
	switch v := item.(type) {
	case types.ByteArray:
		return v
	}
	return nil
}