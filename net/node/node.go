package node

import (
	"dad-go/common"
	"dad-go/common/log"
	. "dad-go/config"
	"dad-go/core/ledger"
	"dad-go/core/transaction"
	. "dad-go/net/message"
	. "dad-go/net/protocol"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"time"
)

// The node capability flag
const (
	RELAY        = 0x01
	SERVER       = 0x02
	NODESERVICES = 0x01
)

type node struct {
	state          uint      // node status
	id             uint64    // The nodes's id
	cap            uint32    // The node capability set
	version        uint32    // The network protocol the node used
	services       uint64    // The services the node supplied
	relay          bool      // The relay capability of the node (merge into capbility flag)
	height         uint64    // The node latest block height
	// TODO does this channel should be a buffer channel
	chF   chan func() error // Channel used to operate the node without lock
	link			// The link status and infomation
	local  *node		// The pointer to local node
	nbrNodes		// The neighbor node connect with currently node except itself
	eventQueue                // The event queue to notice notice other modules
	TXNPool                   // Unconfirmed transaction pool
	idCache                   // The buffer to store the id of the items which already be processed
	ledger     *ledger.Ledger // The Local ledger
}

func (node node) DumpInfo() {
	fmt.Printf("Node info:\n")
	fmt.Printf("\t state = %d\n", node.state)
	fmt.Printf("\t id = 0x%x\n", node.id)
	fmt.Printf("\t addr = %s\n", node.addr)
	fmt.Printf("\t conn = %v\n", node.conn)
	fmt.Printf("\t cap = %d\n", node.cap)
	fmt.Printf("\t version = %d\n", node.version)
	fmt.Printf("\t services = %d\n", node.services)
	fmt.Printf("\t port = %d\n", node.port)
	fmt.Printf("\t relay = %v\n", node.relay)
	fmt.Printf("\t height = %v\n", node.height)

	fmt.Printf("\t conn cnt = %v\n", node.link.connCnt)
}

func (node *node) UpdateInfo(t time.Time, version uint32, services uint64,
	port uint16, nonce uint64, relay uint8, height uint64) {
	// TODO need lock
	node.UpdateTime(t)
	node.id = nonce
	node.version = version
	node.services = services
	node.port = port
	if relay == 0 {
		node.relay = false
	} else {
		node.relay = true
	}
	node.height = uint64(height)
}

func NewNode() *node {
	n := node{
		state: INIT,
		chF:   make(chan func() error),
	}
	runtime.SetFinalizer(&n, rmNode)
	go n.backend()
	return &n
}

func InitNode() Tmper {
	var err error
	n := NewNode()

	n.version = PROTOCOLVERSION
	n.services = NODESERVICES
	n.link.port = uint16(Parameters.NodePort)
	n.relay = true
	rand.Seed(time.Now().UTC().UnixNano())
	// Fixme replace with the real random number
	n.id = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	fmt.Printf("Init node ID to 0x%0x \n", n.id)
	n.nbrNodes.init()
	n.local = n
	n.TXNPool.init()
	n.eventQueue.init()
	n.ledger, err = ledger.GetDefaultLedger()
	if err != nil {
		fmt.Printf("Get Default Ledger error\n")
		errors.New("Get Default Ledger error")
	}

	go n.initConnection()
	go n.updateNodeInfo()
	return n
}

func rmNode(node *node) {
	fmt.Printf("Remove node %s\n", node.addr)
}

// TODO pass pointer to method only need modify it
func (node *node) backend() {
	for f := range node.chF {
		f()
	}
}

func (node node) GetID() uint64 {
	return node.id
}

func (node node) GetState() uint {
	return node.state
}

func (node node) getConn() net.Conn {
	return node.conn
}

func (node node) GetPort() uint16 {
	return node.port
}

func (node node) GetRelay() bool {
	return node.relay
}

func (node node) Version() uint32 {
	return node.version
}

func (node node) Services() uint64 {
	return node.services
}

func (node *node) SetState(state uint) {
	node.state = state
}

func (node *node) LocalNode() Noder {
	return node.local
}

func (node node) GetHeight() uint64 {
	return node.height
}

func (node node) GetLedger() *ledger.Ledger {
	return node.ledger
}

func (node *node) UpdateTime(t time.Time) {
	node.time = t
}

func (node node) GetMemoryPool() map[common.Uint256]*transaction.Transaction {
	return node.GetTxnPool()
	// TODO refresh the pending transaction pool
}

func (node node) SynchronizeMemoryPool() {
	// Fixme need lock
	for _, n := range node.nbrNodes.List {
		if n.state == ESTABLISH {
			ReqMemoryPool(n)
		}
	}
}

func (node node) Xmit(inv common.Inventory) error {

	fmt.Println("****** node Xmit ********")
	var buffer []byte
	var err error

	if inv.Type() == common.TRANSACTION {
		fmt.Printf("****TX transaction message*****\n")
		transaction, isTransaction := inv.(*transaction.Transaction)
		if isTransaction {
			//transaction.Serialize(tmpBuffer)
			buffer, err = NewTx(transaction)
			if err != nil {
				fmt.Println("Error New Tx message ", err.Error())
				return err
			}
		}

	} else if inv.Type() == common.BLOCK {
		fmt.Printf("****TX block message****\n")
		block, isBlock := inv.(*ledger.Block)
		if isBlock {
			buffer, err = NewBlock(block)
			if err != nil {
				fmt.Println("Error New Block message ", err.Error())
				return err
			}
		}
	} else if inv.Type() == common.CONSENSUS {
		fmt.Printf("*****TX consensus message****\n")
		payload, isConsensusPayload := inv.(*ConsensusPayload)
		if isConsensusPayload {
			buffer, err = NewConsensus(payload)
			if err != nil {
				fmt.Println("Error New consensus message ", err.Error())
				return err
			}
		}
	}  else {
		log.Info("Unknow Xmit message type")
		return errors.New("Unknow Xmit message type\n")
 	}

	node.nbrNodes.Broadcast(buffer)

	return nil
}

func (node node) GetAddr() string {
	return node.addr
}

func (node node) GetAddr16() ([16]byte, error) {
	var result [16]byte
	ip := net.ParseIP(node.addr).To16()
	if ip == nil {
		fmt.Printf("Parse IP address error\n")
		return result, errors.New("Parse IP address error")
	}

	copy(result[:], ip[:16])
	return result, nil
}

func (node node) GetTime() int64 {
	t := time.Now()
	return t.UnixNano()
}

func (node node) GetNeighborAddrs() ([]NodeAddr, uint64) {
	var i uint64
	var addrs []NodeAddr
	// TODO read lock
	for _, n := range node.nbrNodes.List {
		if n.GetState() != ESTABLISH {
			continue
		}
		var addr NodeAddr
		addr.IpAddr, _ = n.GetAddr16()
		addr.Time = n.GetTime()
		addr.Services = n.Services()
		addr.Port = n.GetPort()
		addr.ID = n.GetID()
		addrs = append(addrs, addr)

		i++
	}

	return addrs, i
}
