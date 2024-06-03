package raft

import (
	"errors"
	"fmt"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/client"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/app"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

type RaftNode struct {
	address              *Address
	nodeType             NodeType
	log                  []string
	app                  app.KVStore
	electionTerm         int
	clusterAddrList      []Address
	clusterLeaderAddress Address
	electionTimeout      time.Duration
	client               *client.GRPCClient
}

func NewRaftNode(app *app.KVStore, address *Address, isContact bool, contactAddress *Address) *RaftNode {
	raft := &RaftNode{
		app:             *app,
		address:         address,
		nodeType:        LEADER,
		electionTerm:    0,
		clusterAddrList: make([]Address, 0),
		electionTimeout: RandomElectionTimeout(4, 5),
	}

	if isContact {
		
	}

	return raft
}

func (raft RaftNode) initAsLeader() {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) leaderHeartbeat() {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) tryToApplyMembership(contact Address) {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) sendRequest(request interface{}, rpcName string, address Address) {
	fmt.Println("Method Not Implemented")
}

func (raft RaftNode) Heartbeat(request interface{}) interface{} {
	return errors.New("Method Not Implemented")
}

func (raft RaftNode) Execute(request interface{}) interface{} {
	return errors.New("Method Not Implemented")
}
