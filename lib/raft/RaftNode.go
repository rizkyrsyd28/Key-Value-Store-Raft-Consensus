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
	CommittedLength      uint64
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
		ElectionTerm:       TERM_INITIAL,
		CommittedLength:    0,
		VotedFor:           nil,
		ClusterAddressList: ClusterNodeList{Map: map[string]ClusterNode{}},
		ElectionTimeout:    RandomElectionTimeout(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX),
		HeartbeatInterval:  time.Duration(HEARTBEAT_INTERVAL) * time.Second,
		Client:             _client,
		UncommitMembership: nil,
	}

	raft.initPersistentStorage()
	if !isContact {
		logger.InfoLogger.Println("Initalize as Leader")
		raft.initAsLeader()
	} else {
		logger.InfoLogger.Println("Try to apply membership to", contactAddress.ToString(), "from", address.ToString())
		raft.tryToApplyMembership(contactAddress)
	}

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
	persistentStorage := raft.PersistentStorage.Load()
	if persistentStorage != nil {
		fmt.Println("Persistent storage loaded: ", persistentStorage)

		raft.log = persistentStorage.Log
		raft.VotedFor = persistentStorage.VotedFor
		raft.ElectionTerm = uint32(persistentStorage.ElectionTerm)
		raft.CommittedLength = persistentStorage.CommittedLength

		logEntries := persistentStorage.Log.Entries
		for _, entry := range logEntries {
			// Execute log entry
			if entry != nil {
				fmt.Println("Executing log entry: ", entry)
				raft.App.Execute(entry.Command)
			}
		}
		return
	} else {
		raft.PersistentStorage.InitialStoreAll()
		fmt.Println("Persistent storage created")
	}

}

func (raft RaftNode) ResetElectionTimer() {
	raft.timer.Reset(raft.ElectionTimeout)
}

func (raft RaftNode) ResetHeartbeatTimer() {
	raft.timer.Reset(raft.HeartbeatInterval)
}

func (raft *RaftNode) startNode() {
	// fmt.Println(raft.NodeType)
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
	persistentVars := raft.PersistentStorage.Load()
	persistentVars.ElectionTerm += 1
	persistentVars.VotedFor = raft.Address

	logger.DebugLogger.Println("Now persvar: ", persistentVars)
	raft.PersistentStorage.StoreAll(persistentVars)

	responseVote := make(chan pb.RequestVoteResponse)
	votedCount := 1

	fmt.Println("Requesting Vote")
	var wait sync.WaitGroup
	for _, contact := range contactList {
		contact := contact
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
					logger.InfoLogger.Println(err.Error())
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

func (raft *RaftNode) tryToApplyMembership(contact *Address) {
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

	raft.ClusterAddressList.SetAddressPb(response.ClusterAddressList)
	fmt.Println(contact.ToString())
	raft.ClusterLeaderAddress = NewAddress(contact.IP, fmt.Sprint(contact.Port))
	fmt.Println(raft.ClusterLeaderAddress.ToString())
	if response.Log != nil && len(response.Log.Entries) > 0 {
		raft.log.RaftNodeLog = response.Log
		stableValues := raft.PersistentStorage.Load()
		stableValues.Log = raft.log
		raft.PersistentStorage.StoreAll(stableValues)
	}
}

func (raft RaftNode) sendRequest(request interface{}, rpcName string, address Address) {
	fmt.Println("Method Not Implemented")
}

func (raft *RaftNode) sendHeartbeat() {
	type HeartbeatResponseWithAddress struct {
		Status        pb.STATUS
		Term          uint64
		Ack           uint32
		SuccessAppend bool
		Address       Address
	}
	contactList := raft.ClusterAddressList.GetAllAddress()
	responseChan := make(chan HeartbeatResponseWithAddress)
	var wait sync.WaitGroup
	persistentVars := raft.PersistentStorage.Load()
	raft.ElectionTerm = uint32(persistentVars.ElectionTerm)
	raft.CommittedLength = persistentVars.CommittedLength
	raft.log = persistentVars.Log
	raft.VotedFor = persistentVars.VotedFor
	for _, contact := range contactList {
		if contact.IsEqual(raft.Address) {
			continue
		}
		wait.Add(1)
		go func(addr Address) {
			defer wait.Done()
			raft.Client.SetAddress(&addr)
			clusterNode := raft.ClusterAddressList.Get(addr.ToString())
			prefixLen := clusterNode.SentLn
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
				logger.InfoLogger.Println("Error While Send Heartbeat to", addr.ToString())
			} else {
				logger.InfoLogger.Println("sending heartbeat to: ", addr.ToString())
				responseChan <- HeartbeatResponseWithAddress{
					Status:        res.Status,
					Term:          res.Term,
					Ack:           res.Ack,
					SuccessAppend: res.SuccessAppend,
					Address:       addr,
				}
			}
		}(contact)
	}
	// Close the channel after all goroutines finish sending responses
	go func() {
		wait.Wait()
		close(responseChan)
	}()
	for response := range responseChan {
		// TO DO: handle if the response is not as expected
		if response.Status != pb.STATUS_SUCCESS {
			fmt.Printf("Request from %v:%v not succeed - %v\n", response.Address.IP, response.Address.Port, response.Status)
			return
		}
		fmt.Printf("Received response from %v:%v - %v\n", response.Address.IP, response.Address.Port, response.Status)
		respTerm := response.Term
		ack := response.Ack
		successAppend := response.SuccessAppend
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
			raft.ElectionTerm = uint32(respTerm)
			raft.PersistentStorage.StoreAll(persistentVars)
			raft.NodeType = FOLLOWER
			return
		}
	}
}

// func (raft *RaftNode) sendHeartbeat() {
// 	type HeartbeatResponseWithAddress struct {
// 		Response pb.HeartbeatResponse
// 		Address  Address
// 	}

// 	contactList := raft.ClusterAddressList.GetAllAddress()

// 	// responseChan := make(chan pb.HeartbeatResponse)
// 	responseChan := make(chan HeartbeatResponseWithAddress)
// 	var wait sync.WaitGroup
// 	persistentVars := raft.PersistentStorage.Load()

// 	for _, contact := range contactList {
// 		wait.Add(1)
// 		go func(address Address) {
// 			defer wait.Done()
// 			raft.Client.SetAddress(&contact)

// 			clusterNode := raft.ClusterAddressList.Get(contact.ToString())
// 			prefixLen := clusterNode.SentLn
// 			fmt.Println("prefLn", prefixLen)
// 			fmt.Println("stabVa", persistentVars)
// 			fmt.Println("contact list: ", contactList)
// 			suffix := persistentVars.Log.Entries[prefixLen:]
// 			var prefixTerm uint64
// 			prefixTerm = 0
// 			if prefixLen > 0 {
// 				prefixTerm = persistentVars.Log.Entries[prefixLen-1].Term
// 			}

// 			res, err := raft.Client.Services.Raft.SendHeartbeat(context.Background(), &pb.HeartbeatRequest{
// 				Sender:         raft.Address.Address,
// 				Term:           uint64(raft.ElectionTerm),
// 				PrefixLength:   uint64(prefixLen),
// 				PrefixTerm:     prefixTerm,
// 				Suffix:         suffix,
// 				CommitLength:   uint64(persistentVars.CommitLength),
// 				ClusterAddress: raft.ClusterAddressList.GetAllPbAddress(),
// 			})
// 			if err != nil {
// 				fmt.Println("Error While Send Heartbeat")
// 			} else {
// 				fmt.Println("sending heartbeat to: ", contact)
// 				// responseChan <- *res
// 				responseChan <- HeartbeatResponseWithAddress{
// 					Response: *res,
// 					Address:  contact,
// 				}

// 			}
// 		}(contact)
// 	}
// 	// Close the channel after all goroutines finish sending responses
// 	go func() {
// 		wait.Wait()
// 		close(responseChan)
// 	}()
// 	// Now you can range over the responseChan to receive responses
// 	for response := range responseChan {
// 		// TO DO: handle if the response is not as expected

// 		if response.Response.Status != pb.STATUS_SUCCESS {
// 			fmt.Printf("Request from %v:%v not succeed - %v\n", response.Address.IP, response.Address.Port, response.Response.Status)

// 			return
// 		}

// 		fmt.Printf("Received response from %v:%v - %v\n", response.Address.IP, response.Address.Port, response.Response.Status)

// 		respTerm := response.Response.Term
// 		ack := response.Response.Ack
// 		successAppend := response.Response.SuccessAppend

// 		clusterNode := raft.ClusterAddressList.Get(response.Address.ToString())
// 		ackedLen := clusterNode.AckLn

// 		if respTerm == persistentVars.ElectionTerm && raft.NodeType == LEADER {
// 			if successAppend && ack >= uint32(ackedLen) {
// 				clusterNode.SentLn = int64(ack)
// 				clusterNode.AckLn = int64(ack)
// 				raft.ClusterAddressList.PatchAddress(&response.Address, clusterNode)

// 				// TODO: implement commitLogEntries
// 				// raft.commitLogEntries(persistentVars)
// 			} else if clusterNode.SentLn > 0 {
// 				clusterNode.SentLn = clusterNode.SentLn - 1
// 				raft.ClusterAddressList.PatchAddress(&response.Address, clusterNode)
// 			}
// 		} else if respTerm > persistentVars.ElectionTerm {
// 			persistentVars.ElectionTerm = respTerm
// 			raft.PersistentStorage.StoreAll(persistentVars)
// 			raft.NodeType = FOLLOWER
// 			return
// 		}
// 	}
// }

func (raft *RaftNode) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	persistentVars := raft.PersistentStorage.Load()
	raft.VotedFor = persistentVars.VotedFor
	raft.ElectionTerm = uint32(persistentVars.ElectionTerm)
	raft.log = persistentVars.Log
	raft.CommittedLength = persistentVars.CommittedLength

	// leaderAddr := req.LeaderAddress
	reqTerm := req.Term
	prefixLen := int(req.PrefixLength)
	prefixTerm := req.PrefixTerm
	leaderCommit := int(req.CommitLength)
	suffix := req.Suffix
	clusterAddrs := req.ClusterAddress

	if reqTerm > uint64(persistentVars.ElectionTerm) {
		persistentVars.ElectionTerm = uint64(reqTerm)
		persistentVars.VotedFor = nil
		raft.PersistentStorage.StoreAll(persistentVars)
	}

	if reqTerm == uint64(persistentVars.ElectionTerm) {
		raft.NodeType = FOLLOWER
		// raft.ClusterLeaderAddress = &Address{Address: leaderAddr}
		raft.ClusterLeaderAddress = &Address{Address: &pb.Address{IP: req.Sender.IP, Port: req.Sender.Port}}
		// raft.ClusterLeaderAddress = NewAddress(req.LeaderAddress.IP, fmt.Sprint(req.LeaderAddress.Port))
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
	if len(log) == 0 || log[0] == nil {
		return
	}

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
			Status:          pb.STATUS_REDIRECTED.Enum(),
		}, err
	}
	// Check if idempotent then give immediate result
	typeCommand := app.GetCommandType(command)
	if typeCommand == "" {
		errorMessage := "Invalid Command"
		return &pb.Response{
			ResponseMsg: &errorMessage,
			Status:      pb.STATUS_FAILED.Enum(),
		}, nil

	}
	if app.IsCommandIdempotent(typeCommand) {
		value, err := raft.App.Execute(command)
		if err != nil {
			errorMessage := err.Error()
			return &pb.Response{
				ResponseMsg: &errorMessage,
				Status:      pb.STATUS_FAILED.Enum(),
			}, nil
		}
		return &pb.Response{
			Value:  &value,
			Status: pb.STATUS_SUCCESS.Enum(),
		}, nil
	}
	// Append to newLogEntry only, no execute on app

	newLogEntry := &pb.RaftLogEntry{
		Term:    uint64(raft.ElectionTerm),
		Command: command,
	}
	raft.log.Entries = append(raft.log.Entries, newLogEntry)
	raft.PersistentStorage.StoreAll(&persistent_storage.PersistValues{
		ElectionTerm:    uint64(raft.ElectionTerm),
		VotedFor:        raft.VotedFor,
		Log:             raft.log,
		CommittedLength: raft.CommittedLength,
	})
	clusterNode := raft.ClusterAddressList.Get(raft.Address.ToString())
	clusterNode.AckLn = int64(len(raft.log.Entries))
	raft.ClusterAddressList.PatchAddress(raft.Address, clusterNode)
	responseMessageOnProcess := "Process On Progress"
	return &pb.Response{
		ResponseMsg: &responseMessageOnProcess,
		Status:      pb.STATUS_ONPROCESS.Enum(),
	}, nil
}

func (raft *RaftNode) AddMembership(address Address, insert bool) {
	raft.UncommitMembership = NewMembershipApply(address, insert)
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

func (raft *RaftNode) RequestRaftLog() *pb.Response {
	return &pb.Response{
		Log: raft.log.RaftNodeLog,
	}
}
