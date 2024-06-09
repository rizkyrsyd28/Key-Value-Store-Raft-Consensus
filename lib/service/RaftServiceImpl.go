package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/raft"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

type RaftServiceImpl struct {
	pb.UnimplementedRaftServiceServer
	node *raft.RaftNode
}

func NewRaftService(_node *raft.RaftNode) *RaftServiceImpl {
	return &RaftServiceImpl{node: _node}
}

func (rs *RaftServiceImpl) ApplyMembership(ctx context.Context, request *pb.ApplyMembershipRequest) (*pb.ApplyMembershipResponse, error) {
	if rs.node.NodeType != LEADER {
		return &pb.ApplyMembershipResponse{
			Log:                nil,
			ClusterAddressList: nil,
			Status:             pb.STATUS_REDIRECTED,
			RedirectedAddress:  rs.node.ClusterLeaderAddress.Address,
		}, nil
	}

	if !request.Insert && rs.node.ClusterLeaderAddress.IsEqual(&Address{Address: request.Sender}) {
		return &pb.ApplyMembershipResponse{
			Status:            pb.STATUS_FAILED,
			RedirectedAddress: rs.node.ClusterLeaderAddress.Address,
		}, nil
	}

	if !rs.node.NodeMutex.TryLock() {
		message := "Node In Process Of Other Membership Apply"
		return &pb.ApplyMembershipResponse{
			Status:      pb.STATUS_FAILED,
			ResponseMsg: &message,
		}, nil
	}

	defer rs.node.NodeMutex.Unlock()

	applicant := Address{Address: request.Sender}
	contactAddress := append([]Address(nil), rs.node.ClusterAddressList.GetAllAddress()...)

	var wait sync.WaitGroup
	contactResponse := make(chan bool, len(contactAddress))

	internalCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, contact := range contactAddress {
		wait.Add(1)
		go func(_address Address) {
			defer wait.Done()
			rs.node.Client.SetAddress(&_address)
			_, err := rs.node.Client.Services.Raft.AddUpdateCluster(internalCtx, &pb.AddUpdateClusterRequest{
				Insert: request.Insert,
				Sender: rs.node.Address.Address,
			})
			if err != nil {
				contactResponse <- false
			} else {
				contactResponse <- true
			}
		}(contact)
	}

	go func() {
		wait.Wait()
		close(contactResponse)
	}()

	for val := range contactResponse {
		fmt.Println("After: ", val)
	}

	memebershipAcceptCount := 0
	if len(rs.node.ClusterAddressList.Map) > 1 {
		for accept := range contactResponse {
			if accept {
				memebershipAcceptCount++
				if memebershipAcceptCount+1 > len(rs.node.ClusterAddressList.Map)/2 {
					break
				}
			}
		}
	}

	for _, contact := range contactAddress {
		go func(contact Address) { // Salin variabel contact ke dalam goroutine
			rs.node.Client.SetAddress(&contact)
			rs.node.Client.Services.Raft.CommitUpdateCluster(internalCtx, &pb.CommitUpdateClusterRequest{
				Insert: request.Insert,
				Sender: rs.node.Address.Address,
			})
		}(contact) // Kirim salinan variabel contact ke goroutine
	}

	// fmt.Println("Data : ", rs.node.UncommitMembership)

	if request.Insert {
		rs.node.ClusterAddressList.AddAddress(&applicant)
	}

	fmt.Println("Interrupt And Restart Loop")

	contact := []*pb.Address(nil)
	for _, value := range rs.node.ClusterAddressList.Map {
		contact = append(contact, value.Address.Address)
	}

	response := &pb.ApplyMembershipResponse{
		Status:             pb.STATUS_SUCCESS,
		Log:                rs.node.RequestRaftLog().GetLog(),
		ClusterAddressList: contact,
	}

	return response, nil
}

func (rs *RaftServiceImpl) RequestVote(ctx context.Context, request *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {

	// TODO: add voted for check in condition
	if rs.node.ElectionTerm == request.Term && rs.node.NodeType == FOLLOWER {
		return &pb.RequestVoteResponse{
			NodeAddress: rs.node.Address.Address,
			CurrentTerm: 1,
			Granted:     true,
		}, nil
	} else {
		return &pb.RequestVoteResponse{
			NodeAddress: rs.node.Address.Address,
			CurrentTerm: 1,
			Granted:     false,
		}, nil
	}
}

func (rs *RaftServiceImpl) SendHeartbeat(ctx context.Context, request *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	rs.node.ResetElectionTimer()
	rs.node.ClusterLeaderAddress = &Address{Address: request.Sender}
	response, err := rs.node.Heartbeat(ctx, request)

	if err != nil {
		return nil, err
	}

	// response := pb.HeartbeatResponse{
	// 	Ack:           35,
	// 	Term:          1,
	// 	SuccessAppend: true,
	// }

	return response, nil
}

func (rs *RaftServiceImpl) AddUpdateCluster(ctx context.Context, request *pb.AddUpdateClusterRequest) (*pb.AddUpdateClusterResponse, error) {
	rs.node.AddMembership(Address{Address: request.Sender}, request.Insert)

	response := pb.AddUpdateClusterResponse{Status: pb.STATUS_SUCCESS}

	return &response, nil
}

func (rs *RaftServiceImpl) CommitUpdateCluster(ctx context.Context, request *pb.CommitUpdateClusterRequest) (*pb.CommitUpdateClusterResponse, error) {

	if err := rs.node.CommitMembership(Address{Address: request.Sender}, request.Insert); err != nil {
		message := fmt.Sprintf("%v", err.Error())
		return &pb.CommitUpdateClusterResponse{
			Status:      pb.STATUS_FAILED,
			ResponseMsg: &message,
		}, err
	}

	return &pb.CommitUpdateClusterResponse{
		Status: pb.STATUS_SUCCESS,
	}, nil
}
