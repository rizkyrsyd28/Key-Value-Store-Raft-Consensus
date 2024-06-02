package lib

import (
	"errors"
	"fmt"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/enum"
	_struct "github.com/Sister20/if3230-tubes-dark-syster/lib/struct"
)

type RaftNode struct {
	address              _struct.Address
	nodeType             enum.NodeType
	log                  []string
	app                  KVStore
	electionTerm         int
	clusterAddrList      []_struct.Address
	clusterLeaderAddress _struct.Address
}

func New(app *KVStore, address *_struct.Address, contactAddress _struct.Address) *RaftNode {
	raft := &RaftNode{
		app:             *app,
		address:         *address,
		nodeType:        enum.LEADER,
		electionTerm:    0,
		clusterAddrList: make([]_struct.Address, 0),
	}

	// if contactAddress == nil {

	// }

	return raft
}

func (raft RaftNode) initAsLeader() {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) leaderHeartbeat() {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) tryToApplyMembership(contact _struct.Address) {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) sendRequest(request interface{}, rpcName string, address _struct.Address) {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) Heartbeat(request interface{}) interface{} {
	return errors.New("Method Not Implemented")
}

func (raft RaftNode) Execute(request interface{}) interface{} {
	return errors.New("Method Not Implemented")
}
