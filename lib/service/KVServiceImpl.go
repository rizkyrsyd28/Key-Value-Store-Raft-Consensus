package service

import (
	"context"
	"strconv"

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

func (kv *KeyValueServiceImpl) Ping(context.Context, *pb.Empty) (*pb.Response, error) {
	res := kv.node.App.Ping()
	response := pb.Response{Value: res}
	return &response, nil
}

func (kv *KeyValueServiceImpl) Get(ctx context.Context, req *pb.KeyRequest) (*pb.Response, error) {
	res, err := kv.node.App.Get(req.Key)
	if err != nil {
		return nil, err
	}
	response := pb.Response{Value: res}
	return &response, nil
}

func (kv *KeyValueServiceImpl) Set(ctx context.Context, req *pb.KeyValueRequest) (*pb.Response, error) {
	res, err := kv.node.App.Set(req.Key, req.Value)
	if err != nil {
		return nil, err
	}
	response := pb.Response{Value: res}
	return &response, nil
}
func (kv *KeyValueServiceImpl) StrLn(ctx context.Context, req *pb.KeyRequest) (*pb.Response, error) {
	res, err := kv.node.App.Strlen(req.Key)
	if err != nil {
		return nil, err
	}
	response := pb.Response{Value: strconv.Itoa(res)}
	return &response, nil
}
func (kv *KeyValueServiceImpl) Del(ctx context.Context, req *pb.KeyRequest) (*pb.Response, error) {
	res, err := kv.node.App.Delete(req.Key)
	if err != nil {
		return nil, err
	}
	response := pb.Response{Value: res}
	return &response, nil
}
func (kv *KeyValueServiceImpl) Append(ctx context.Context, req *pb.KeyValueRequest) (*pb.Response, error) {
	res, err := kv.node.App.Append(req.Key, req.Value)
	if err != nil {
		return nil, err
	}
	response := pb.Response{Value: res}
	return &response, nil
}
