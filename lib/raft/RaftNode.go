package raft

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/client"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/stable_storage"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/util"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/app"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/logger"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

type RaftNode struct {
	Address              *Address
	NodeType             NodeType
	log                  logger.RaftNodeLog
	App                  app.KVStore
	StableStorage        *stable_storage.StableStorage
	ElectionTerm         int
	ClusterAddressList   ClusterNodeList
	ClusterLeaderAddress *Address
	VotesReceived        []Address
	ElectionTimeout      time.Duration
	Client               *client.GRPCClient
	NodeMutex            sync.Mutex // goroutine
	UncommitMembership   *MembershipApply

	interruptChan chan struct{}
	shutdownChan  chan struct{}
	wg            sync.WaitGroup
}

func NewRaftNode(app *app.KVStore, address *Address, isContact bool, contactAddress *Address) *RaftNode {

	_client, err := client.NewClient(address)
	if err != nil {
		fmt.Println("Error")
	}

	raft := &RaftNode{
		App:                *app,
		Address:            address,
		NodeType:           FOLLOWER,
		log:                logger.RaftNodeLog{},
		StableStorage:      &stable_storage.StableStorage{},
		ElectionTerm:       0,
		ClusterAddressList: ClusterNodeList{Map: map[string]ClusterNode{}},
		ElectionTimeout:    RandomElectionTimeout(4, 5),
		Client:             _client,
		UncommitMembership: nil,
	}

	if !isContact {
		raft.ClusterAddressList.AddAddress(address)
		raft.initAsLeader()
	} else {
		raft.tryToApplyMembership(contactAddress)
	}

	return raft
}

func (raft *RaftNode) initAsLeader() {
	// Print Log Initalize as Leader
	raft.ClusterLeaderAddress = raft.Address
	raft.NodeType = LEADER
}

// func (raft *RaftNode) leaderHeartbeat() {
// 	if raft.NodeType == LEADER {
// 		addrs := raft.ClusterAddressList.GetAllAddress()
// 		fmt.Printf("Sending heartbeat... %v\n", addrs)

// 		for _, addr := range addrs {
// 			raft.ThreadPool.Submit(raft.sendHeartbeat, addr)
// 		}
// 	}

// 	raft.interruptAndRestartLoop()
// }

func (raft *RaftNode) leaderHeartbeat() {
	if raft.NodeType == LEADER {
		addrs := raft.ClusterAddressList.GetAllAddress()
		fmt.Printf("Sending heartbeat... %v\n", addrs)

		var wg sync.WaitGroup
		for _, addr := range addrs {
			wg.Add(1)
			go func(addr *Address) {
				defer wg.Done()
				raft.sendHeartbeat(addr.ToString())
			}(&addr)
		}
		wg.Wait()
	}

	raft.interruptAndRestartLoop()
}

func (raft *RaftNode) sendHeartbeat(id string) {
	clusterNode := raft.ClusterAddressList.Get(id)
	addr := clusterNode.Address

	if addr == raft.Address {
		return
	}

	stableVars := raft.StableStorage.Load()
	prefixLen := clusterNode.SentLn
	suffix := stableVars.Log.Entries[prefixLen:]
	var prefixTerm uint64
	prefixTerm = 0
	if prefixLen > 0 {
		prefixTerm = stableVars.Log.Entries[prefixLen-1].Term
	}

	msg := &pb.HeartbeatRequest{
		LeaderAddress:  raft.Address.Address,
		Term:           uint64(raft.ElectionTerm),
		PrefixLength:   uint64(prefixLen),
		PrefixTerm:     prefixTerm,
		Suffix:         suffix,
		CommitLength:   uint64(stableVars.CommitLength),
		ClusterAddress: raft.ClusterAddressList.GetAllPbAddress(),
	}

	// resp, err := raft.RPCHandler.Request(addr, string("Heartbeat"), msg)
	resp, err := raft.Client.Services.Raft.SendHeartbeat(context.Background(), msg)
	if err != nil {
		return
	}

	if resp.Status != pb.STATUS_SUCCESS {
		return
	}

	respTerm := resp.Term
	ack := resp.Ack
	successAppend := resp.SuccessAppend

	ackedLen := clusterNode.AckLn

	stableVars = raft.StableStorage.Load()
	if respTerm == stableVars.ElectionTerm && raft.NodeType == LEADER {
		if successAppend && ack >= uint32(ackedLen) {
			clusterNode.SentLn = int64(ack)
			clusterNode.AckLn = int64(ack)
			raft.ClusterAddressList.PatchAddress(addr, clusterNode)

			// TODO: implement commitLogEntries
			// raft.commitLogEntries(stableVars)
		} else if clusterNode.SentLn > 0 {
			clusterNode.SentLn = clusterNode.SentLn - 1
			raft.ClusterAddressList.PatchAddress(addr, clusterNode)

			// TODO: REPLICATE LOG
		}
	} else if respTerm > stableVars.ElectionTerm {
		stableVars.ElectionTerm = respTerm
		stableVars.VotedFor = nil
		raft.StableStorage.StoreAll(stableVars)
		raft.NodeType = FOLLOWER
		raft.VotesReceived = make([]Address, 0)
		raft.interruptAndRestartLoop()
		return
	}
}

func (raft RaftNode) tryToApplyMembership(contact *Address) {
	if err := raft.Client.SetAddress(contact); err != nil {
		fmt.Println("Error client while change address")
		return
	}

	request := &pb.ApplyMembershipRequest{
		Insert: true,
		Sender: raft.Address.Address,
	}

	response, err := raft.Client.Services.Raft.ApplyMembership(context.Background(), request)
	if err != nil {
		raft.ClusterAddressList.AddAddress(raft.Address)
		fmt.Printf("Error While Apply %v\n", err.Error())
		return
	}

	for _, data := range response.ClusterAddressList {
		addr := Address{Address: data}
		fmt.Println(addr.ToString())
	}

	raft.ClusterAddressList.SetAddressPb(response.ClusterAddressList)
	raft.ClusterLeaderAddress = contact
}

func (raft RaftNode) sendRequest(request interface{}, rpcName string, address Address) {
	fmt.Println("Method Not Implemented")
}

func (raft *RaftNode) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	leaderAddr := req.LeaderAddress
	reqTerm := req.Term
	prefixLen := int(req.PrefixLength)
	prefixTerm := req.PrefixTerm
	leaderCommit := int(req.CommitLength)
	suffix := req.Suffix
	clusterAddrs := req.ClusterAddress

	stableVars := raft.StableStorage.Load()
	if reqTerm > uint64(stableVars.ElectionTerm) {
		stableVars.ElectionTerm = uint64(reqTerm)
		stableVars.VotedFor = nil
		raft.StableStorage.StoreAll(stableVars)
		raft.interruptAndRestartLoop()
	}

	if reqTerm == uint64(stableVars.ElectionTerm) {
		raft.NodeType = FOLLOWER
		raft.VotesReceived = make([]Address, 0)
		raft.ClusterLeaderAddress = &Address{Address: leaderAddr}
	}

	log := stableVars.Log.Entries
	logOk := len(log) >= prefixLen && (prefixLen == 0 || log[prefixLen-1].Term == prefixTerm)

	response := &pb.HeartbeatResponse{
		Status: pb.STATUS_SUCCESS,
		Term:   uint64(stableVars.ElectionTerm),
	}

	if reqTerm == uint64(stableVars.ElectionTerm) && logOk {
		raft.ClusterAddressList.SetAddressPb(clusterAddrs)
		raft.appendEntries(prefixLen, leaderCommit, suffix, stableVars)
		ack := prefixLen + len(suffix)
		response.Ack = uint32(ack)
		response.SuccessAppend = true
		raft.interruptAndRestartLoop()
	} else {
		response.Ack = 0
		response.SuccessAppend = false
	}

	return response, nil
}

func (raft *RaftNode) appendEntries(prefixLen int, leaderCommit int, suffix []*pb.RaftLogEntry, stableVars *stable_storage.StableVars) {
	log := stableVars.Log.Entries

	if len(suffix) > 0 && len(log) > prefixLen {
		idx := FindMin(len(log), prefixLen+len(suffix)) - 1
		if log[idx].Term != suffix[idx-prefixLen].Term {
			log = log[:prefixLen]
		}
	}

	if prefixLen+len(suffix) > len(log) {
		for i := len(log) - prefixLen; i < len(suffix); i++ {
			log = append(log, suffix[i])
		}
	}

	stableVars.Log.Entries = log

	commitLength := stableVars.CommitLength
	if uint64(leaderCommit) > commitLength {
		for i := commitLength; i < uint64(leaderCommit); i++ {
			raft.App.Execute(log[i].Command)
		}
		stableVars.CommitLength = uint64(leaderCommit)
	}

	raft.StableStorage.StoreAll(stableVars)
}

func (raft RaftNode) Execute(request interface{}) interface{} {
	return errors.New("Method Not Implemented")
}

func (raft *RaftNode) AddMembership(address Address, insert bool) {
	raft.UncommitMembership = NewMembershipApply(address, insert)

	fmt.Println("Data\n", raft.UncommitMembership)
}

func (raft *RaftNode) CommitMembership(address Address, insert bool) error {

	if raft.UncommitMembership == nil {
		return errors.New("Uncommited Member Not Found")
	}

	if !raft.UncommitMembership.Address.IsEqual(&address) || raft.UncommitMembership.Insert != insert {
		return errors.New("Uncommited Member Not Match")
	}

	if insert {
		raft.ClusterAddressList.AddAddress(&address)
	}

	return nil
}

func (raft *RaftNode) startLoop() {
	raft.interruptChan = make(chan struct{})
	raft.shutdownChan = make(chan struct{})

	raft.wg.Add(1)
	go raft.runLoop()
}

func (raft *RaftNode) runLoop() {
	defer raft.wg.Done()

	ticker := time.NewTicker(raft.getLoopInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			raft.performLoopTask()
		case <-raft.interruptChan:
			return
		case <-raft.shutdownChan:
			return
		}
	}
}

const HeartbeatInterval = time.Duration(util.HEARTBEAT_INTERVAL) * time.Second

func (raft *RaftNode) getLoopInterval() time.Duration {
	if raft.NodeType == LEADER {
		return HeartbeatInterval
	}

	// Interval waktu random untuk pemilihan
	return time.Duration(RandomElectionTimeout(int(ELECTION_TIMEOUT_MIN), int(ELECTION_TIMEOUT_MAX)).Nanoseconds())
}

func (raft *RaftNode) performLoopTask() {
	if raft.NodeType == LEADER {
		raft.leaderHeartbeat()
	} else {
		// TODO: requestVote or other tasks
		// raft.requestVote()
		fmt.Println("Request Vote")
	}
}

func (raft *RaftNode) interruptAndRestartLoop() {
	raft.interruptChan <- struct{}{}
}

func (raft *RaftNode) stopLoop() {
	close(raft.shutdownChan)
	raft.wg.Wait()
}
