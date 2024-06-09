package raft

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/client"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/persistent_storage"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/app"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/logger"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

type RaftNode struct {
	Address              *Address
	NodeType             NodeType
	log                  logger.RaftNodeLog
	ElectionTerm         uint32
	App                  *app.KVStore
	PersistentStorage    *persistent_storage.PersistentStorage
	CommittedLength         uint64
	VotedFor             *Address
	ClusterAddressList   ClusterNodeList
	ClusterLeaderAddress *Address
	ElectionTimeout      time.Duration // in seconds
	HeartbeatInterval    time.Duration // in seconds
	Client               *client.GRPCClient
	NodeMutex            sync.Mutex // goroutine
	UncommitMembership   *MembershipApply
	timer                *time.Timer
}

func NewRaftNode(app *app.KVStore, address *Address, isContact bool, contactAddress *Address) *RaftNode {

	_client, err := client.NewClient(address)
	if err != nil {
		fmt.Println("Error")
	}

	raft := &RaftNode{
		App:      app,
		Address:  address,
		NodeType: FOLLOWER,
		log: logger.RaftNodeLog{
			RaftNodeLog: &pb.RaftNodeLog{
				Entries: make([]*pb.RaftLogEntry, 0),
			},
		},
		PersistentStorage:  persistent_storage.NewPersistentStorage(address),
		ElectionTerm:       0,
		CommittedLength:       0,
		VotedFor:           nil,
		ClusterAddressList: ClusterNodeList{Map: map[string]ClusterNode{}},
		ElectionTimeout:    RandomElectionTimeout(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX),
		HeartbeatInterval:  time.Duration(HEARTBEAT_INTERVAL) * time.Second,
		Client:             _client,
		UncommitMembership: nil,
	}

	if !isContact {
		raft.initAsLeader()
	} else {
		raft.tryToApplyMembership(contactAddress)
	}

	raft.initPersistentStorage()
	raft.startNode()

	return raft
}

func (raft *RaftNode) initAsLeader() {
	// Print Log Initalize as Leader
	raft.ClusterLeaderAddress = raft.Address
	raft.ClusterAddressList.AddAddress(raft.Address)
	raft.NodeType = LEADER
	fmt.Println("Leader Initalize")
}

func (raft *RaftNode) initPersistentStorage() {
	raft.PersistentStorage = persistent_storage.NewPersistentStorage(raft.Address)
	res := raft.PersistentStorage.Load()
	if res != nil {
		fmt.Println("Persistent storage loaded: ", res)

		logEntries := res.Log.Entries
		for _, entry := range logEntries {
			// Execute log entry
			if entry != nil {
				fmt.Println("Executing log entry: ", entry)
				raft.App.Execute(entry.Command)
			}
		}
		return
	} else {
		RaftLog := logger.RaftNodeLog{
			RaftNodeLog: &pb.RaftNodeLog{
				Entries: make([]*pb.RaftLogEntry, 1),
			},
		}
		dummy := persistent_storage.PersistValues{
			ElectionTerm: 1,
			VotedFor:     raft.Address,
			Log:          RaftLog,
			CommittedLength: 0,
		}

		raft.PersistentStorage.StoreAll(&dummy)
		fmt.Println("Persistent storage created: ", res)
	}

}

func (raft RaftNode) ResetElectionTimer() {
	raft.timer.Reset(raft.ElectionTimeout)
}

func (raft RaftNode) ResetHeartbeatTimer() {
	raft.timer.Reset(raft.HeartbeatInterval)
}

func (raft *RaftNode) startNode() {
	fmt.Println(raft.NodeType)
	if raft.NodeType == LEADER {
		raft.timer = time.NewTimer(raft.HeartbeatInterval)
	} else {
		raft.timer = time.NewTimer(raft.ElectionTimeout)
	}
	go func() {
		for {
			select {
			case <-raft.timer.C:
				if raft.NodeType == LEADER {
					raft.sendHeartbeat()
					raft.ResetHeartbeatTimer()
				} else {
					raft.requestVote() // reset election timer inside requestVote
				}
			}
			fmt.Println("=====")
		}
	}()
}

func (raft *RaftNode) requestVote() {
	raft.NodeType = CANDIDATE

	contactList := raft.ClusterAddressList.GetAllAddress()

	// TODO: update election term on stable storage

	responseVote := make(chan pb.RequestVoteResponse)
	votedCount := 1

	fmt.Println("Requesting Vote")
	var wait sync.WaitGroup
	for _, contact := range contactList {
		wait.Add(1)
		go func(address Address) {
			defer wait.Done() // send request vote to other nodes
			if contact.Address != raft.Address.Address {
				raft.Client.SetAddress(&contact)
				res, err := raft.Client.Services.Raft.RequestVote(context.Background(), &pb.RequestVoteRequest{
					VotedFor: raft.Address.Address,
					Term:     raft.ElectionTerm,
					// LogLength: persistentStorage.length,
					// LogTerm: persistentStorage.term,
					// Sender: raft.Address.Address,
				})
				if err != nil {
					fmt.Println("Error While Send Heartbeat")
				} else {
					fmt.Println("sending heartbeat...")
					responseVote <- *res
				}
			}
		}(contact)
	}

	// rerandom election timeout and reset timer
	raft.ElectionTimeout = RandomElectionTimeout(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX)
	raft.timer.Reset(raft.ElectionTimeout)

	go func() {
		wait.Wait()
		close(responseVote)
	}()

	for res := range responseVote {
		if res.Granted {
			votedCount++
		}
	}

	// Check if votedCount is enough to become the leader
	if votedCount > len(contactList)/2 {
		raft.initAsLeader()
		raft.ResetHeartbeatTimer()
		fmt.Println("Became the leader, votedCount:", votedCount)
	} else {
		fmt.Println("Failed to become the leader, votedCount:", votedCount)
		// timer continue running
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

func (raft *RaftNode) sendHeartbeat() {
	type HeartbeatResponseWithAddress struct {
		Response pb.HeartbeatResponse
		Address  Address
	}

	contactList := raft.ClusterAddressList.GetAllAddress()

	// responseChan := make(chan pb.HeartbeatResponse)
	responseChan := make(chan HeartbeatResponseWithAddress)
	var wait sync.WaitGroup
	persistentVars := raft.PersistentStorage.Load()

	for _, contact := range contactList {
		wait.Add(1)
		go func(address Address) {
			defer wait.Done()
			raft.Client.SetAddress(&contact)

			clusterNode := raft.ClusterAddressList.Get(contact.ToString())
			prefixLen := clusterNode.SentLn
			fmt.Println("prefLn", prefixLen)
			fmt.Println("stabVa", persistentVars)
			suffix := persistentVars.Log.Entries[prefixLen:]
			var prefixTerm uint64
			prefixTerm = 0
			if prefixLen > 0 {
				prefixTerm = persistentVars.Log.Entries[prefixLen-1].Term
			}

			res, err := raft.Client.Services.Raft.SendHeartbeat(context.Background(), &pb.HeartbeatRequest{
				Sender:         raft.Address.Address,
				Term:           uint64(raft.ElectionTerm),
				PrefixLength:   uint64(prefixLen),
				PrefixTerm:     prefixTerm,
				Suffix:         suffix,
				CommitLength:   uint64(persistentVars.CommittedLength),
				ClusterAddress: raft.ClusterAddressList.GetAllPbAddress(),
			})
			if err != nil {
				fmt.Println("Error While Send Heartbeat")
			} else {
				fmt.Println("sending heartbeat...")
				// responseChan <- *res
				responseChan <- HeartbeatResponseWithAddress{
					Response: *res,
					Address:  contact,
				}

			}
		}(contact)
	}
	// Close the channel after all goroutines finish sending responses
	go func() {
		wait.Wait()
		close(responseChan)
	}()
	// Now you can range over the responseChan to receive responses
	for response := range responseChan {
		// TO DO: handle if the response is not as expected

		if response.Response.Status != pb.STATUS_SUCCESS {
			fmt.Printf("Request from %v:%v not succeed - %v\n", response.Address.IP, response.Address.Port, response.Response.Status)

			return
		}

		fmt.Printf("Received response from %v:%v - %v\n", response.Address.IP, response.Address.Port, response.Response.Status)

		respTerm := response.Response.Term
		ack := response.Response.Ack
		successAppend := response.Response.SuccessAppend

		clusterNode := raft.ClusterAddressList.Get(response.Address.ToString())
		ackedLen := clusterNode.AckLn

		if respTerm == persistentVars.ElectionTerm && raft.NodeType == LEADER {
			if successAppend && ack >= uint32(ackedLen) {
				clusterNode.SentLn = int64(ack)
				clusterNode.AckLn = int64(ack)
				raft.ClusterAddressList.PatchAddress(&response.Address, clusterNode)

				// TODO: implement commitLogEntries
				raft.executeLogEntries(persistentVars)
			} else if clusterNode.SentLn > 0 {
				clusterNode.SentLn = clusterNode.SentLn - 1
				raft.ClusterAddressList.PatchAddress(&response.Address, clusterNode)
			}
		} else if respTerm > persistentVars.ElectionTerm {
			persistentVars.ElectionTerm = respTerm
			raft.PersistentStorage.StoreAll(persistentVars)
			raft.NodeType = FOLLOWER
			return
		}
	}
}

func (raft *RaftNode) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	leaderAddr := req.LeaderAddress
	reqTerm := req.Term
	prefixLen := int(req.PrefixLength)
	prefixTerm := req.PrefixTerm
	leaderCommit := int(req.CommitLength)
	suffix := req.Suffix
	clusterAddrs := req.ClusterAddress

	persistentVars := raft.PersistentStorage.Load()
	if reqTerm > uint64(persistentVars.ElectionTerm) {
		persistentVars.ElectionTerm = uint64(reqTerm)
		persistentVars.VotedFor = nil
		raft.PersistentStorage.StoreAll(persistentVars)
	}

	if reqTerm == uint64(persistentVars.ElectionTerm) {
		raft.NodeType = FOLLOWER
		raft.ClusterLeaderAddress = &Address{Address: leaderAddr}
	}

	log := persistentVars.Log.Entries
	logOk := len(log) >= prefixLen && (prefixLen == 0 || log[prefixLen-1].Term == prefixTerm)

	response := &pb.HeartbeatResponse{
		Status: pb.STATUS_SUCCESS,
		Term:   uint64(persistentVars.ElectionTerm),
	}

	if reqTerm == uint64(persistentVars.ElectionTerm) && logOk {
		raft.ClusterAddressList.SetAddressPb(clusterAddrs)
		raft.appendEntries(prefixLen, leaderCommit, suffix, persistentVars)
		ack := prefixLen + len(suffix)
		response.Ack = uint32(ack)
		response.SuccessAppend = true
	} else {
		response.Ack = 0
		response.SuccessAppend = false
	}

	return response, nil
}

func (raft *RaftNode) appendEntries(prefixLen int, leaderCommit int, suffix []*pb.RaftLogEntry, persistentVars *persistent_storage.PersistValues) {
	log := persistentVars.Log.Entries

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

	persistentVars.Log.Entries = log

	commitLength := persistentVars.CommittedLength
	if uint64(leaderCommit) > commitLength {
		for i := commitLength; i < uint64(leaderCommit); i++ {
			raft.App.Execute(log[i].Command)
		}
		persistentVars.CommittedLength = uint64(leaderCommit)
	}

	raft.PersistentStorage.StoreAll(persistentVars)
}

func (raft *RaftNode) executeLogEntries(persistentValues *persistent_storage.PersistValues) {
	// 50% + 1 mechanism
	minAck := math.Ceil(float64(len(raft.ClusterAddressList.GetAllAddress())+1) / 2)
	log := persistentValues.Log.Entries
	if len(log) == 0 || log[0] == nil {
		return
	}
	latestIdxToCommitFromAck := 0
	for i := 0; i < len(log); i++ {
		sum := 0
		for _, clusterNode := range raft.ClusterAddressList.Map {
			if clusterNode.AckLn >= int64(i) {
				sum++
			}
		}
		if sum >= int(minAck) {
			latestIdxToCommitFromAck = i + 1
		}
	}
	commitLength := persistentValues.CommittedLength
	latestLogElectionTerm := log[latestIdxToCommitFromAck-1].Term
	currerntElectionTerm := persistentValues.ElectionTerm
	if latestIdxToCommitFromAck > int(commitLength) && latestLogElectionTerm == currerntElectionTerm {
		for i := commitLength; i < uint64(latestIdxToCommitFromAck); i++ {
			raft.App.Execute(log[i].Command)
		}
		persistentValues.CommittedLength = uint64(latestIdxToCommitFromAck)
		raft.PersistentStorage.StoreAll(persistentValues)
	}
}

func (raft *RaftNode) Execute(ctx context.Context, command string) (*pb.Response, error) {
	if raft.NodeType != LEADER {
		err := errors.New("redirected to leader address")
		return &pb.Response{
			RedirectAddress: raft.ClusterLeaderAddress.Address,
		}, err
	}
	// Check if idempotent then give immediate result
	typeCommand := app.GetCommandType(command)
	if typeCommand == "" {
		errorMessage := "Invalid Command"
		return &pb.Response{
			ResponseMsg: &errorMessage,
		}, nil

	}
	if app.IsCommandIdempotent(typeCommand) {
		value, err := raft.App.Execute(command)
		if err != nil {
			errorMessage := err.Error()
			return &pb.Response{
				ResponseMsg: &errorMessage,
			}, nil
		}
		return &pb.Response{
			Value: &value,
		}, nil
	}
	// Append to newLogEntry only, no execute on app
	newLogEntry := &pb.RaftLogEntry{
		Term:    uint64(raft.ElectionTerm),
		Command: command,
	}
	raft.log.Entries = append(raft.log.Entries, newLogEntry)
	raft.PersistentStorage.StoreAll(&persistent_storage.PersistValues{
		ElectionTerm: uint64(raft.ElectionTerm),
		VotedFor:     raft.VotedFor,
		Log:          raft.log,
		CommittedLength: raft.CommittedLength,
	})
	clusterNode := raft.ClusterAddressList.Get(raft.Address.ToString())
	clusterNode.AckLn = int64(len(raft.log.Entries))
	raft.ClusterAddressList.PatchAddress(raft.Address, clusterNode)
	responseMessageOnProcess := "Process On Progress"
	return &pb.Response{
		ResponseMsg: &responseMessageOnProcess,
	}, nil
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
