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

package exec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/dad-go/common"
	"github.com/dad-go/vm/neovm/interfaces"
	"github.com/dad-go/vm/wasmvm/memory"
	"github.com/dad-go/vm/wasmvm/validate"
	"github.com/dad-go/vm/wasmvm/wasm"
	"math"
	"os"
	"reflect"
)

const (
	CONTRACT_METHOD_NAME = "invoke"
	PARAM_SPLITER        = "|"
)

//todo add parameters
func NewExecutionEngine(icontainer interfaces.ICodeContainer, icrypto interfaces.ICrypto, itable interfaces.ICodeTable, iservice IInteropService, ver string) *ExecutionEngine {

	engine := &ExecutionEngine{
		crypto: icrypto, table: itable,
		codeContainer: icontainer,
		service:       NewInteropService(),
		version:       ver,
	}
	if iservice != nil {
		engine.service.MergeMap(iservice.GetServiceMap())
	}
	return engine
}

type ExecutionEngine struct {
	crypto        interfaces.ICrypto
	table         interfaces.ICodeTable
	service       *InteropService
	codeContainer interfaces.ICodeContainer
	//memory  		*memory.VMmemory
	vm      *VM
	version string //for test different contracts
}

func (e *ExecutionEngine) GetVM() *VM {
	return e.vm
}

//todo use this method just for test
func (e *ExecutionEngine) CallInf(caller common.Address, code []byte, input []interface{}, message []interface{}) ([]byte, error) {
	methodad-gome := input[0].(string)

	//1. read code
	bf := bytes.NewBuffer(code)

	//2. read module
	m, err := wasm.ReadModule(bf, importer)
	if err != nil {
		return nil, errors.New("Verify wasm failed!" + err.Error())
	}

	//3. verify the module
	//already verified in step 2

	//4. check the export
	//every wasm should have at least 1 export
	if m.Export == nil {
		return nil, errors.New("No export in wasm!")
	}

	vm, err := NewVM(m)
	if err != nil {
		return nil, err
	}
	vm.Engine = e
	if e.service != nil {
		vm.Services = e.service.GetServiceMap()
	}
	e.vm = vm
	vm.Engine = e

	vm.SetMessage(message)

	vm.Caller = caller
	vm.CodeHash = common.ToCodeHash(code)

	entry, ok := m.Export.Entries[methodad-gome]
	if ok == false {
		return nil, errors.New("Method:" + methodad-gome + " does not exist!")
	}
	//get entry index
	index := int64(entry.Index)

	//get function index
	fidx := m.Function.Types[int(index)]

	//get  function type
	ftype := m.Types.Entries[int(fidx)]

	paramlength := len(input) - 1
	if len(ftype.ParamTypes) != paramlength {
		return nil, errors.New("parameter count is not right")
	}
	params := make([]uint64, paramlength)
	for i, param := range input[1:] {
		//if type is struct
		if reflect.TypeOf(param).Kind() == reflect.Struct {
			offset, err := vm.SetStructMemory(param)
			if err != nil {
				return nil, err
			}
			params[i] = uint64(offset)
		} else {
			switch param.(type) {
			case string:
				offset, err := vm.SetPointerMemory(param.(string))
				if err != nil {
					return nil, err
				}
				params[i] = uint64(offset)

				/*				offset, err := vm.SetMemory(param)
								if err != nil {
									return nil, err
								}
								vm.GetMemory().MemPoints[uint64(offset)] = &memory.TypeLength{Ptype:memory.P_STRING,Length:len(param.(string))}
								params[i] = uint64(offset)*/
			case int:
				params[i] = uint64(param.(int))
			case int64:
				params[i] = uint64(param.(int64))
			case float32:
				bits := math.Float32bits(param.(float32))
				params[i] = uint64(bits)
			case float64:
				bits := math.Float64bits(param.(float64))
				params[i] = uint64(bits)

			case []int:
				idx := 0
				for i, v := range param.([]int) {
					offset, err := vm.SetMemory(v)
					if err != nil {
						return nil, err
					}
					if i == 0 {
						idx = offset
					}
				}
				vm.GetMemory().MemPoints[uint64(idx)] = &memory.TypeLength{Ptype: memory.P_INT32, Length: len(param.([]int)) * 4}
				params[i] = uint64(idx)

			case []int64:
				idx := 0
				for i, v := range param.([]int64) {
					offset, err := vm.SetMemory(v)
					if err != nil {
						return nil, err
					}
					if i == 0 {
						idx = offset
					}
				}
				vm.GetMemory().MemPoints[uint64(idx)] = &memory.TypeLength{Ptype: memory.P_INT64, Length: len(param.([]int64)) * 8}
				params[i] = uint64(idx)

			case []float32:
				idx := 0
				for i, v := range param.([]float32) {
					offset, err := vm.SetMemory(v)
					if err != nil {
						return nil, err
					}
					if i == 0 {
						idx = offset
					}
				}
				vm.GetMemory().MemPoints[uint64(idx)] = &memory.TypeLength{Ptype: memory.P_FLOAT32, Length: len(param.([]float32)) * 4}
				params[i] = uint64(idx)
			case []float64:
				idx := 0
				for i, v := range param.([]float64) {
					offset, err := vm.SetMemory(v)
					if err != nil {
						return nil, err
					}
					if i == 0 {
						idx = offset
					}
				}
				vm.GetMemory().MemPoints[uint64(idx)] = &memory.TypeLength{Ptype: memory.P_FLOAT64, Length: len(param.([]float64)) * 8}
				params[i] = uint64(idx)
			}
		}

	}

	res, err := vm.ExecCode(false, index, params...)
	if err != nil {
		return nil, errors.New("ExecCode error!" + err.Error())
	}

	if len(ftype.ReturnTypes) == 0 {
		return nil, nil
	}

	switch ftype.ReturnTypes[0] {
	case wasm.ValueTypeI32:
		return Int32ToBytes(res.(uint32)), nil
	case wasm.ValueTypeI64:
		return Int64ToBytes(res.(uint64)), nil
	case wasm.ValueTypeF32:
		return Float32ToBytes(res.(float32)), nil
	case wasm.ValueTypeF64:
		return Float64ToBytes(res.(float64)), nil
	default:
		return nil, errors.New("the return type is not supported")
	}

	return nil, nil
}

func (e *ExecutionEngine) GetMemory() *memory.VMmemory {
	return e.vm.memory
}

//func (e *ExecutionEngine)GetMemStruct() *memory.VMmemory{
//	return e.memory
//}

func (e *ExecutionEngine) Create(caller common.Address, code []byte) ([]byte, error) {
	return code, nil
}

//the input format should be "methodad-gome | args"
func (e *ExecutionEngine) Call(caller common.Address, code, input []byte) (returnbytes []byte, er error) {

	//catch the panic to avoid crash the whole node
	defer func() {
		if err := recover(); err != nil {
			returnbytes = nil
			er = errors.New("error happend")
		}
	}()

	if e.version != "test" {

		methodad-gome := CONTRACT_METHOD_NAME //fix to "invoke"

		tmparr := bytes.Split(input, []byte(PARAM_SPLITER))
		if len(tmparr) != 2 {
			return nil, errors.New("input format is not right!")
		}

		//1. read code
		bf := bytes.NewBuffer(code)

		//2. read module
		m, err := wasm.ReadModule(bf, importer)
		if err != nil {
			return nil, errors.New("Verify wasm failed!" + err.Error())
		}

		//3. verify the module
		//already verified in step 2

		//4. check the export
		//every wasm should have at least 1 export
		if m.Export == nil {
			return nil, errors.New("No export in wasm!")
		}

		vm, err := NewVM(m)
		if err != nil {
			return nil, err
		}
		if e.service != nil {
			vm.Services = e.service.GetServiceMap()
		}
		e.vm = vm
		vm.Engine = e
		//todo add message from input
		//vm.SetMessage(message)

		vm.Caller = caller
		vm.CodeHash = common.ToCodeHash(code)

		entry, ok := m.Export.Entries[methodad-gome]
		if ok == false {
			return nil, errors.New("Method:" + methodad-gome + " does not exist!")
		}
		//get entry index
		index := int64(entry.Index)

		//get function index
		fidx := m.Function.Types[int(index)]

		//get  function type
		ftype := m.Types.Entries[int(fidx)]

		params := make([]uint64, 2)

		actionName := string(tmparr[0])
		actIdx, err := vm.SetPointerMemory(actionName)
		if err != nil {
			return nil, err
		}
		params[0] = uint64(actIdx)

		args := tmparr[1]
		argIdx, err := vm.SetPointerMemory(args)
		if err != nil {
			return nil, err
		}
		//init paramIndex
		vm.memory.ParamIndex = 0

		params[1] = uint64(argIdx)

		res, err := vm.ExecCode(false, index, params...)
		if err != nil {
			return nil, errors.New("ExecCode error!" + err.Error())
		}

		if len(ftype.ReturnTypes) == 0 {
			return nil, nil
		}

		switch ftype.ReturnTypes[0] {
		case wasm.ValueTypeI32:
			return Int32ToBytes(res.(uint32)), nil
		case wasm.ValueTypeI64:
			return Int64ToBytes(res.(uint64)), nil
		case wasm.ValueTypeF32:
			return Float32ToBytes(res.(float32)), nil
		case wasm.ValueTypeF64:
			return Float64ToBytes(res.(float64)), nil
		default:
			return nil, errors.New("the return type is not supported")
		}

	} else {
		//for test
		methodad-gome, err := getCallMethodad-gome(input)
		if err != nil {
			return nil, err
		}

		//1. read code
		bf := bytes.NewBuffer(code)

		//2. read module
		m, err := wasm.ReadModule(bf, importer)
		if err != nil {
			return nil, errors.New("Verify wasm failed!" + err.Error())
		}

		//3. verify the module
		//already verified in step 2

		//4. check the export
		//every wasm should have at least 1 export
		if m.Export == nil {
			return nil, errors.New("No export in wasm!")
		}

		vm, err := NewVM(m)
		if err != nil {
			return nil, err
		}
		if e.service != nil {
			vm.Services = e.service.GetServiceMap()
		}
		e.vm = vm
		vm.Engine = e
		//todo add message from input
		//vm.SetMessage(message)
		entry, ok := m.Export.Entries[methodad-gome]
		if ok == false {
			return nil, errors.New("Method:" + methodad-gome + " does not exist!")
		}
		//get entry index
		index := int64(entry.Index)

		//get function index
		fidx := m.Function.Types[int(index)]

		//get  function type
		ftype := m.Types.Entries[int(fidx)]

		//paramtypes := ftype.ParamTypes

		params, err := getParams(input)
		if err != nil {
			return nil, err
		}

		if len(params) != len(ftype.ParamTypes) {
			return nil, errors.New("Parameters count is not right")
		}

		res, err := vm.ExecCode(false, index, params...)
		if err != nil {
			return nil, errors.New("ExecCode error!" + err.Error())
		}

		if len(ftype.ReturnTypes) == 0 {
			return nil, nil
		}

		switch ftype.ReturnTypes[0] {
		case wasm.ValueTypeI32:
			return Int32ToBytes(res.(uint32)), nil
		case wasm.ValueTypeI64:
			return Int64ToBytes(res.(uint64)), nil
		case wasm.ValueTypeF32:
			return Float32ToBytes(res.(float32)), nil
		case wasm.ValueTypeF64:
			return Float64ToBytes(res.(float64)), nil
		default:
			return nil, errors.New("the return type is not supported")
		}

	}

}

//TODO NOT IN USE BUT DON'T DELETE IT
//current we only support the ONT SYSTEM module import
//other imports will raise an error
func importer(name string) (*wasm.Module, error) {
	//TODO add the path into config file
	if name != "ONT" {
		return nil, errors.New("import [" + name + "] is not supported! ")
	}
	f, err := os.Open(name + ".wasm")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, err := wasm.ReadModule(f, nil)
	err = validate.VerifyModule(m)
	if err != nil {
		return nil, err
	}
	return m, nil

}

//get call method name from the input bytes
//the input should be:[Namelength][methodad-gome][paramcount][param1Length][param2Length].....[param1Data][Param2Data][....]
//input[0] should be the name length
//next n bytes should the be the method name
func getCallMethodad-gome(input []byte) (string, error) {

	if len(input) <= 1 {
		return "", errors.New("input format error!")
	}

	length := int(input[0])

	if length > len(input[1:]) {
		return "", errors.New("input method name length error!")
	}

	return string(input[1 : length+1]), nil
}

func getParams(input []byte) ([]uint64, error) {
	//log.Error(fmt.Sprintf("in getParams: input is %v\n",input))

	nameLength := int(input[0])

	paramCnt := int(input[1+nameLength])

	res := make([]uint64, paramCnt)

	paramlengthSlice := input[1+nameLength+1 : 1+1+nameLength+paramCnt]

	paramSlice := input[1+nameLength+1+paramCnt:]

	for i := 0; i < paramCnt; i++ {
		//get param length
		pl := int(paramlengthSlice[i])

		if (i+1)*pl > len(paramSlice) {
			return nil, errors.New("get param failed!")
		}
		param := paramSlice[i*pl : (i+1)*pl]

		if len(param) < 8 {
			temp := make([]byte, 8)
			copy(temp, param)
			res[i] = binary.LittleEndian.Uint64(temp)
		} else {
			res[i] = binary.LittleEndian.Uint64(param)
		}
	}
	//log.Error(res)
	return res, nil

}

func Int32ToBytes(i32 uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, i32)
	return bytesBuffer.Bytes()
}

func Int64ToBytes(i64 uint64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, i64)
	return bytesBuffer.Bytes()
}
func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
