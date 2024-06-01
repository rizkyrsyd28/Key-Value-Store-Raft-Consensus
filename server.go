package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

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

func main() {
	netListen, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err.Error())
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor))
	ExecuteService := lib.ExecuteServiceImpl{}
	pb.RegisterExecuteServiceServer(server, &ExecuteService)

	log.Println("Server Started At ", netListen.Addr())

	if err := server.Serve(netListen); err != nil {
		log.Fatal(err.Error())
	}
}
