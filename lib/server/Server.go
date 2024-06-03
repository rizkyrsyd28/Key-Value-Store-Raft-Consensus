package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/raft"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/app"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type GRPCServer struct {
	address    *Address
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(_address *Address, isContact bool, contactAddress *Address) *GRPCServer {
	server := &GRPCServer{
		address: _address,
	}

	netListen, err := net.Listen("tcp", server.address.ToString())
	if err != nil {
		log.Fatal(err.Error())
	}
	server.listener = netListen

	server.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor))

	app := app.NewKVStore()
	raft := raft.NewRaftNode(app, server.address, isContact, contactAddress)

	kvservice := service.NewKVService(raft)
	raftservice := service.NewRaftService(raft)
	pb.RegisterKeyValueServiceServer(server.grpcServer, kvservice)
	pb.RegisterRaftServiceServer(server.grpcServer, raftservice)

	return server
}

func (server *GRPCServer) Serve() {
	if err := server.grpcServer.Serve(server.listener); err != nil {
		log.Fatal(err.Error())
	}
}

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Log method name and request
	p, _ := peer.FromContext(ctx)
	log.Printf("Method: %s, Request: %+v, From: %v", info.FullMethod, req, p.Addr)

	// Handling request
	resp, err := handler(ctx, req)

	// Log response and elapsed time
	log.Printf("Method: %s, Response: %+v, ElapsedTime: %s, Error: %v", info.FullMethod, resp, time.Since(start), err)

	return resp, err
}
