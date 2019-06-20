package ledger

import (
	"dad-go/common"
	"dad-go/core/contract/program"
	"io"
	"dad-go/common/serialization"
)


type Blockdata struct {
	//TODO: define the Blockheader struct(define new uinttype)
	Version uint
	PrevBlockHash  common.Uint256
	TransactionsRoot common.Uint256
	Timestamp uint
	Height uint
	consensusData uint64
	Program *program.Program

	hash common.Uint256
}

//Serialize the blockheader
func (bd *Blockdata) Serialize(w io.Writer)  {
	bd.SerializeUnsigned(w)
	w.Write([]byte{byte(1)})
	bd.Program.Serialize(w)
}

//Serialize the blockheader data without program
func (bd *Blockdata) SerializeUnsigned(w io.Writer) error  {
	//TODO: implement blockheader SerializeUnsigned

	return nil
}

func (bd *Blockdata) Deserialize(r io.Reader)  {
	//TODO：Blockdata Deserialize
}

func (bd *Blockdata) DeserializeUnsigned(r io.Reader)  error{
	bd.Version = serialization.ReadUint(r)

	//prevBlock
	preBlock := new(common.Uint256)
	err := preBlock.Deserialize(r)
	if err != nil {
		return err
	}
	bd.PrevBlockHash = *preBlock

	//TransactionsRoot
	txRoot := new(common.Uint256)
	err = txRoot.Deserialize(r)
	if err != nil {	return err}
	bd.TransactionsRoot = *txRoot

	//Timestamp
	bd.Timestamp = serialization.ReadUint(r)

	//Height
	bd.Height = serialization.ReadUint(r)

	//consensusData
	bd.consensusData = serialization.ReadUint64(r)

	return nil
}

func (bd *Blockdata) GetProgramHashes() ([]common.Uint160, error){
	//TODO: implement blockheader GetProgramHashes

	return nil, nil
}



func (bd *Blockdata) SetPrograms(programs []*program.Program){
	if(len(programs) != 1){
		return
	}
	bd.Program = programs[0]
}

func (bd *Blockdata) GetPrograms() []*program.Program {
	return []*program.Program {bd.Program}
}


func (bd *Blockdata) Hash() common.Uint256 {
	//TODO: implement Blockdata Hash

	return common.Uint256{}
}

