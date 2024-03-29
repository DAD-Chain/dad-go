/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */
package ontid

import (
	"bytes"
	"testing"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/smartcontract/service/native"
	"github.com/ontio/dad-go/smartcontract/service/native/testsuite"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
)

func testcase(t *testing.T, f func(t *testing.T, n *native.NativeService)) {
	testsuite.InvokeNativeContract(t, utils.OntIDContractAddress,
		func(n *native.NativeService) ([]byte, error) {
			f(t, n)
			return nil, nil
		},
	)
}

func TestReg(t *testing.T) {
	testcase(t, CaseRegID)
}

func TestOwner(t *testing.T) {
	testcase(t, CaseOwner)
}

func TestOwnerSize(t *testing.T) {
	testcase(t, CaseOwnerSize)
}

// Register id with account acc
func regID(n *native.NativeService, id string, a *account.Account) error {
	// make arguments
	sink := common.NewZeroCopySink(nil)
	sink.WriteVarBytes([]byte(id))
	pk := keypair.SerializePublicKey(a.PubKey())
	sink.WriteVarBytes(pk)
	n.Input = sink.Bytes()
	// set signing address
	n.Tx.SignedAddr = []common.Address{a.Address}
	// call
	_, err := regIdWithPublicKey(n)
	return err
}

func CaseRegID(t *testing.T, n *native.NativeService) {
	id, err := account.GenerateID()
	if err != nil {
		t.Fatal(err)
	}
	a := account.NewAccount("")

	// 1. register invalid id, should fail
	if err := regID(n, "did:ont:abcd1234", a); err == nil {
		t.Error("invalid id registered")
	}

	// 2. register without valid signature, should fail
	sink := common.NewZeroCopySink(nil)
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{}
	if _, err := regIdWithPublicKey(n); err == nil {
		t.Error("id registered without signature")
	}

	// 3. register with invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes([]byte("invalid public key"))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a.Address}
	if _, err := regIdWithPublicKey(n); err == nil {
		t.Error("id registered with invalid key")
	}

	// 4. register id
	if err := regID(n, id, a); err != nil {
		t.Fatal(err)
	}

	// 5. get DDO
	sink.Reset()
	sink.WriteString(id)
	n.Input = sink.Bytes()
	_, err = GetDDO(n)
	if err != nil {
		t.Error(err)
	}

	// 6. register again, should fail
	if err := regID(n, id, a); err == nil {
		t.Error("id registered twice")
	}

	// 7. revoke with invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 2)
	n.Input = sink.Bytes()
	if _, err := revokeID(n); err == nil {
		t.Error("revoked by invalid key")
	}

	// 8. revoke without valid signature, should fail
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 1)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{common.ADDRESS_EMPTY}
	if _, err := revokeID(n); err == nil {
		t.Error("revoked without valid signature")
	}

	// 9. revoke id
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 1)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a.Address}
	if _, err := revokeID(n); err != nil {
		t.Fatal(err)
	}

	// 10. register again, should fail
	if err := regID(n, id, a); err == nil {
		t.Error("revoked id should not be registered again")
	}

	// 11. get DDO of the revoked id
	sink.Reset()
	sink.WriteString(id)
	n.Input = sink.Bytes()
	_, err = GetDDO(n)
	if err == nil {
		t.Error("get DDO of the revoked id should fail")
	}

}

func CaseOwner(t *testing.T, n *native.NativeService) {
	// 1. register ID
	id, err := account.GenerateID()
	if err != nil {
		t.Fatal("generate ID error")
	}
	a0 := account.NewAccount("")
	if err := regID(n, id, a0); err != nil {
		t.Fatal("register ID error", err)
	}

	// 2. add key without valid signature, should fail
	a1 := account.NewAccount("")
	sink := common.NewZeroCopySink(nil)
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{common.ADDRESS_EMPTY}
	if _, err = addKey(n); err == nil {
		t.Error("key added without valid signature")
	}

	// 3. add key by invalid owner, should fail
	a2 := account.NewAccount("")
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a2.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a2.Address}
	if _, err = addKey(n); err == nil {
		t.Error("key added by invalid owner")
	}

	// 4. add invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes([]byte("test invalid key"))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	if _, err = addKey(n); err == nil {
		t.Error("invalid key added")
	}

	// 5. add key
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	if _, err = addKey(n); err != nil {
		t.Fatal(err)
	}

	// 6. verify new key
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 2)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a1.Address}
	res, err := verifySignature(n)
	if err != nil || !bytes.Equal(res, utils.BYTE_TRUE) {
		t.Fatal("verify the added key failed")
	}

	// 7. add key again, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	if _, err = addKey(n); err == nil {
		t.Fatal("should not add the same key twice")
	}

	// 8. remove key without valid signature, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a2.Address}
	if _, err = removeKey(n); err == nil {
		t.Error("key removed without valid signature")
	}

	// 9. remove key by invalid owner, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a2.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a2.Address}
	if _, err = removeKey(n); err == nil {
		t.Error("key removed by invalid owner")
	}

	// 10. remove invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a2.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a1.Address}
	if _, err = removeKey(n); err == nil {
		t.Error("invalid key removed")
	}

	// 11. remove key
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	if _, err = removeKey(n); err != nil {
		t.Fatal(err)
	}

	// 12. check removed key
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 1)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	res, err = verifySignature(n)
	if err == nil && bytes.Equal(res, utils.BYTE_TRUE) {
		t.Fatal("removed key passed verification")
	}

	// 13. add removed key again, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	res, err = verifySignature(n)
	if err == nil && bytes.Equal(res, utils.BYTE_TRUE) {
		t.Error("the removed key should not be added again")
	}

	// 14. query removed key
	sink.Reset()
	sink.WriteString(id)
	sink.WriteInt32(1)
	n.Input = sink.Bytes()
	_, err = GetPublicKeyByID(n)
	if err == nil {
		t.Error("query removed key should fail")
	}
}

func CaseOwnerSize(t *testing.T, n *native.NativeService) {
	id, _ := account.GenerateID()
	a := account.NewAccount("")
	err := regID(n, id, a)
	if err != nil {
		t.Fatal(err)
	}

	enc, err := encodeID([]byte(id))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, OWNER_TOTAL_SIZE)
	_, err = insertPk(n, enc, buf)
	if err == nil {
		t.Fatal("total size of the owner's key should be limited")
	}
}
