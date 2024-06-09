package service

import (
	"context"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/raft"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
)

type KeyValueServiceImpl struct {
	pb.UnimplementedKeyValueServiceServer
	node *raft.RaftNode
}

func NewKVService(_node *raft.RaftNode) *KeyValueServiceImpl {
	return &KeyValueServiceImpl{node: _node}
}

func (kv *KeyValueServiceImpl) Ping(ctx context.Context, empty *pb.Empty) (*pb.Response, error) {
	res, err := kv.node.Execute(ctx, "PING")
	if err != nil {
		return res, err
	}
	return res, nil
}

func (kv *KeyValueServiceImpl) Get(ctx context.Context, req *pb.KeyRequest) (*pb.Response, error) {
	res, err := kv.node.Execute(ctx, ("GET " + req.Key))
	if err != nil {
		return res, err
	}
	return res, nil
}

func (kv *KeyValueServiceImpl) Set(ctx context.Context, req *pb.KeyValueRequest) (*pb.Response, error) {
	res, err := kv.node.Execute(ctx, ("SET " + req.Key + " " + req.Value))
	if err != nil {
		return res, err
	}
	return res, nil
}
func (kv *KeyValueServiceImpl) StrLn(ctx context.Context, req *pb.KeyRequest) (*pb.Response, error) {
	res, err := kv.node.Execute(ctx, ("STRLEN " + req.Key))
	if err != nil {
		return res, err
	}
	return res, nil
}
func (kv *KeyValueServiceImpl) Del(ctx context.Context, req *pb.KeyRequest) (*pb.Response, error) {
	res, err := kv.node.Execute(ctx, ("DELETE " + req.Key))
	if err != nil {
		return res, err
	}
	return res, nil
}
func (kv *KeyValueServiceImpl) Append(ctx context.Context, req *pb.KeyValueRequest) (*pb.Response, error) {
	res, err := kv.node.Execute(ctx, ("APPEND " + req.Key + " " + req.Value))
	if err != nil {
		return res, err
	}
	return res, nil
}
