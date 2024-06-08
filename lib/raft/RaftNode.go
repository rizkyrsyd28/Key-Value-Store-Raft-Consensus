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
	ElectionTerm         int
	ClusterAddressList   ClusterNodeList
	ClusterLeaderAddress *Address
	ElectionTimeout      time.Duration
	HeartbeatInterval     time.Duration
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

	const (
		ELECTION_MIN_TIMEOUT = 15
		ELECTION_MAX_TIMEOUT = 30
		HEARTBEAT_INTERVAL = 4
	)

	raft := &RaftNode{
		App:                *app,
		Address:            address,
		NodeType:           FOLLOWER,
		ElectionTerm:       0,
		ClusterAddressList: ClusterNodeList{Map: map[string]ClusterNode{}},
		ElectionTimeout:    RandomElectionTimeout(ELECTION_MIN_TIMEOUT, ELECTION_MAX_TIMEOUT),
		HeartbeatInterval:   HEARTBEAT_INTERVAL,
		Client:             _client,
		UncommitMembership: nil,
	}

	
	if !isContact {
		raft.ClusterAddressList.AddAddress(address)
		raft.initAsLeader()
	} else {
		raft.tryToApplyMembership(contactAddress)
	}
			
	raft.start()
	

	return raft
}

func (raft *RaftNode) initAsLeader() {
	// Print Log Initalize as Leader
	raft.ClusterLeaderAddress = raft.Address
	fmt.Println("djdjdjjd")
	raft.NodeType = LEADER
}

func (raft RaftNode) ResetTimer() {
	raft.timer.Reset(raft.ElectionTimeout)
}

func (raft *RaftNode) start() {
	fmt.Println("halo ")
	fmt.Println(raft.NodeType)
	if (raft.NodeType == LEADER) {
		fmt.Println("halo ha")
		raft.leaderHeartbeat()
		fmt.Println("halo halo")
	} else {
		raft.runElectionTimer()
	}
}


func (raft RaftNode) runElectionTimer() {
	raft.timer = time.NewTimer(raft.ElectionTimeout * time.Second)
	go func ()  {
		<- raft.timer.C
		raft.requestVote()
	}()
}

func (raft RaftNode) requestVote() {
	raft.NodeType = CANDIDATE
	raft.ElectionTerm++

}

func (raft *RaftNode) leaderHeartbeat() {
	raft.timer = time.NewTimer(time.Duration(raft.HeartbeatInterval) * time.Second)
	go func() {
		for {
			select {
			case <-raft.timer.C:
				raft.sendHeartbeat()
				raft.timer.Reset(time.Duration(raft.HeartbeatInterval) * time.Second)
			}
		}
	}()
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
		// Handle each response here
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
