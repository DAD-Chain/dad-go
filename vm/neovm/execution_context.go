package neovm

import (
	"github.com/dad-go/common"
	"github.com/dad-go/vm/neovm/types"
	"github.com/dad-go/vm/neovm/utils"
	"io"
)

type ExecutionContext struct {
	Code               []byte
	OpReader           *utils.VmReader
	PushOnly           bool
	BreakPoints        []uint
	InstructionPointer int
	CodeHash           common.Uint160
	engine             *ExecutionEngine
}

func NewExecutionContext(engine *ExecutionEngine, code []byte, pushOnly bool, breakPoints []uint) *ExecutionContext {
	var executionContext ExecutionContext
	executionContext.Code = code
	executionContext.OpReader = utils.NewVmReader(code)
	executionContext.PushOnly = pushOnly
	executionContext.BreakPoints = breakPoints
	executionContext.InstructionPointer = 0
	executionContext.engine = engine
	return &executionContext
}

func (ec *ExecutionContext) GetInstructionPointer() int {
	return ec.OpReader.Position()
}

func (ec *ExecutionContext) SetInstructionPointer(offset int64) {
	ec.OpReader.Seek(offset, io.SeekStart)
}

func (ec *ExecutionContext) GetCodeHash() (common.Uint160, error) {
	if ec.CodeHash.CompareTo(common.Uint160{}) == 0 {
		codeHash := common.ToCodeHash(ec.Code)
		ec.CodeHash = codeHash
	}
	return ec.CodeHash, nil
}

func (ec *ExecutionContext) NextInstruction() OpCode {
	return OpCode(ec.Code[ec.OpReader.Position()])
}

func (ec *ExecutionContext) Clone() *ExecutionContext {
	executionContext := NewExecutionContext(ec.engine, ec.Code, ec.PushOnly, ec.BreakPoints)
	executionContext.InstructionPointer = ec.InstructionPointer
	executionContext.SetInstructionPointer(int64(ec.GetInstructionPointer()))
	return executionContext
}

func (ec *ExecutionContext) GetStackItem() types.StackItemInterface {
	return nil
}

func (ec *ExecutionContext) GetExecutionContext() *ExecutionContext {
	return ec
}
