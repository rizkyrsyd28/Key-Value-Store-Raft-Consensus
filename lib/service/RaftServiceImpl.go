package service

import (
	"context"
	"fmt"

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
	fmt.Printf("Request: { %v }\n", request)
	response := pb.ApplyMembershipResponse{
		Log: []string{"A", "B"},
		ClusterAddressList: []*pb.Address{
			&NewAddress("135.0.10.9", "4000").Address,
			&NewAddress("135.0.16.0", "4001").Address,
			&NewAddress("135.0.11.9", "4002").Address,
		},
	}
	return &response, nil
}

func (rs *RaftServiceImpl) RequestVote(ctx context.Context, request *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	address := NewAddress("135.0.12.9", "4005")
	response := pb.RequestVoteResponse{
		NodeAddress: &address.Address,
		CurrentTerm: 1,
		Granted:     true,
	}
	return &response, nil
}

func (rs *RaftServiceImpl) SendHeartbeat(ctx context.Context, request *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	response := pb.HeartbeatResponse{
		Ack:           35,
		Term:          1,
		SuccessAppend: true,
	}
	return &response, nil
}
