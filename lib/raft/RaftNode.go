package raft

import (
	"errors"
	"fmt"
	_struct "github.com/Sister20/if3230-tubes-dark-syster/lib/util"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/app"
)

type RaftNode struct {
	address              *_struct.Address
	nodeType             _struct.NodeType
	log                  []string
	app                  app.KVStore
	electionTerm         int
	clusterAddrList      []_struct.Address
	clusterLeaderAddress _struct.Address
}

func NewRaftNode(app *app.KVStore, address *_struct.Address) *RaftNode {
	raft := &RaftNode{
		app:             *app,
		address:         address,
		nodeType:        _struct.LEADER,
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
