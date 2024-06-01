package lib

import (
	"context"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
)

type ExecuteServiceImpl struct {
	pb.UnimplementedExecuteServiceServer
}

func (e *ExecuteServiceImpl) Ping(context.Context, *pb.Empty) (*pb.Response, error) {
	response := pb.Response{Value: "Pong"}
	return &response, nil
}

func (e *ExecuteServiceImpl) Get(context.Context, *pb.KeyRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Get Data"}
	return &response, nil
}

func (e *ExecuteServiceImpl) Set(context.Context, *pb.KeyValueRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Set Data"}
	return &response, nil
}
func (e *ExecuteServiceImpl) StrLn(context.Context, *pb.KeyRequest) (*pb.Response, error) {
	response := pb.Response{Value: "StrLn Data"}
	return &response, nil
}
func (e *ExecuteServiceImpl) Del(context.Context, *pb.KeyRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Delete Data"}
	return &response, nil
}
func (e *ExecuteServiceImpl) Append(context.Context, *pb.KeyValueRequest) (*pb.Response, error) {
	response := pb.Response{Value: "Append Data"}
	return &response, nil
}
