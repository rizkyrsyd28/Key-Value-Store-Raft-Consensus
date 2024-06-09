package raft

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/client"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/app"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

type RaftNode struct {
	Address              *Address
	NodeType             NodeType
	log                  []string
	App                  app.KVStore
	ElectionTerm         uint32
	ClusterAddressList   ClusterNodeList
	ClusterLeaderAddress *Address
	ElectionTimeout      time.Duration // in seconds
	HeartbeatInterval    time.Duration // in seconds
	Client               *client.GRPCClient
	NodeMutex            sync.Mutex // goroutine
	UncommitMembership   *MembershipApply
	timer 			  	 *time.Timer
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
		ElectionTerm:       0,
		ClusterAddressList: ClusterNodeList{Map: map[string]ClusterNode{}},
		ElectionTimeout:    RandomElectionTimeout(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX),
		HeartbeatInterval:   time.Duration(HEARTBEAT_INTERVAL) * time.Second,
		Client:             _client,
		UncommitMembership: nil,
	}

	
	if !isContact {
		raft.initAsLeader()
	} else {
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

func (raft RaftNode) ResetElectionTimer() {
	raft.timer.Reset(raft.ElectionTimeout)
}

func (raft RaftNode) ResetHeartbeatTimer() {
	raft.timer.Reset(raft.HeartbeatInterval)
}

func (raft *RaftNode) startNode() {
	fmt.Println(raft.NodeType)
	if (raft.NodeType == LEADER) {
		raft.timer = time.NewTimer(raft.HeartbeatInterval)
	} else {
		raft.timer = time.NewTimer(raft.ElectionTimeout)
	}
	go func() {
		for {
			select {
			case <-raft.timer.C:
				if (raft.NodeType == LEADER) {
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
		go func (address Address) {
			defer wait.Done()// send request vote to other nodes
			if (contact.Address != raft.Address.Address) {
				raft.Client.SetAddress(&contact)
				res, err := raft.Client.Services.Raft.RequestVote(context.Background(), &pb.RequestVoteRequest{
					VotedFor: raft.Address.Address,
					Term: raft.ElectionTerm,
					// LogLength: stableStorage.length,
					// LogTerm: stableStorage.term,
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
	contactList := raft.ClusterAddressList.GetAllAddress()
	
	responseChan := make(chan pb.HeartbeatResponse)
	var wait sync.WaitGroup

	for _, contact := range contactList {
		wait.Add(1)
		go func (address Address) {
			defer wait.Done()
			raft.Client.SetAddress(&contact)
			res, err := raft.Client.Services.Raft.SendHeartbeat(context.Background(), &pb.HeartbeatRequest{
				Sender: raft.Address.Address,
			})
			if err != nil {
				fmt.Println("Error While Send Heartbeat")
			} else {
				fmt.Println("sending heartbeat...")
				responseChan <- *res
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

		fmt.Println("Received response:", response)
	}
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
