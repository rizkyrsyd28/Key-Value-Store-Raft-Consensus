package service

import (
	"context"

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
	// if rs.node.NodeType != LEADER {
	// 	return &pb.ApplyMembershipResponse{
	// 		Log:                nil,
	// 		ClusterAddressList: nil,
	// 		Status:             pb.STATUS_REDIRECTED,
	// 		RedirectedAddress:  rs.node.ClusterLeaderAddress.Address,
	// 	}, nil
	// }

	// if !request.Insert && rs.node.ClusterLeaderAddress.IsEqual(&Address{Address: request.Sender}) {
	// 	return &pb.ApplyMembershipResponse{
	// 		Status: pb.STATUS_FAILED,
	// 		RedirectedAddress: rs.node.ClusterLeaderAddress.Address,
	// 	}, nil
	// }

	// return &response, nil
}

func (rs *RaftServiceImpl) RequestVote(ctx context.Context, request *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	address := NewAddress("135.0.12.9", "4005")
	response := pb.RequestVoteResponse{
		NodeAddress: address.Address,
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
