package validation

import (
	"errors"
	. "github.com/dad-go/common"
	sig "github.com/dad-go/core/signature"
	"github.com/dad-go/crypto"
	. "github.com/dad-go/errors"
	vm "github.com/dad-go/vm/neovm"
	"github.com/dad-go/vm/neovm/interfaces"
	"github.com/dad-go/smartcontract/service"
	"github.com/dad-go/smartcontract/types"
)

func VerifySignableData(signableData sig.SignableData) (bool, error) {

	hashes, err := signableData.GetProgramHashes()
	if err != nil {
		return false, err
	}

	programs := signableData.GetPrograms()
	Length := len(hashes)
	if Length != len(programs) {
		return false, errors.New("The number of data hashes is different with number of programs.")
	}

	programs = signableData.GetPrograms()
	for i := 0; i < len(programs); i++ {
		temp, _ := ToCodeHash(programs[i].Code)
		if hashes[i] != temp {
			return false, errors.New("The data hashes is different with corresponding program code.")
		}
		//execute program on VM
		var cryptos interfaces.ICrypto
		cryptos = new(vm.ECDsaCrypto)
		stateReader := service.NewStateReader(types.Verification)
		se := vm.NewExecutionEngine(signableData, cryptos, nil, stateReader)
		se.LoadCode(programs[i].Code, false)
		se.LoadCode(programs[i].Parameter, true)
		se.Execute()

		if se.GetState() != vm.HALT {
			return false, NewDetailErr(errors.New("[VM] Finish State not equal to HALT."), ErrNoCode, "")
		}

		if se.GetEvaluationStack().Count() != 1 {
			return false, NewDetailErr(errors.New("[VM] Execute Engine Stack Count Error."), ErrNoCode, "")
		}

		flag := se.GetExecuteResult()
		if !flag {
			return false, NewDetailErr(errors.New("[VM] Check Sig FALSE."), ErrNoCode, "")
		}
	}

	return true, nil
}

func VerifySignature(signableData sig.SignableData, pubkey *crypto.PubKey, signature []byte) (bool, error) {
	err := crypto.Verify(*pubkey, sig.GetHashData(signableData), signature)
	if err != nil {
		return false, NewDetailErr(err, ErrNoCode, "[Validation], VerifySignature failed.")
	} else {
		return true, nil
	}
}
