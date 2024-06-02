package service

import (
	"context"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
)

type KeyValueServiceImpl struct {
	pb.UnimplementedKeyValueServiceServer
}

func (kv *KeyValueServiceImpl) Ping(context.Context, *pb.Empty) (*pb.Response, error) {
	response := pb.Response{Value: "PONG"}
	return &response, nil
}

func (kv *KeyValueServiceImpl) Get(context.Context, *pb.KeyRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Get Data"}
	return &response, nil
}

func (kv *KeyValueServiceImpl) Set(context.Context, *pb.KeyValueRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Set Data"}
	return &response, nil
}
func (kv *KeyValueServiceImpl) StrLn(context.Context, *pb.KeyRequest) (*pb.Response, error) {
	response := pb.Response{Value: "StrLn Data"}
	return &response, nil
}
func (kv *KeyValueServiceImpl) Del(context.Context, *pb.KeyRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Delete Data"}
	return &response, nil
}
func (kv *KeyValueServiceImpl) Append(context.Context, *pb.KeyValueRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Append Data"}
	return &response, nil
}
