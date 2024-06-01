package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/service"
	_struct "github.com/Sister20/if3230-tubes-dark-syster/lib/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var address _struct.Address

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

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run server.go <server ip> <server port> <time?>")
		return
	}

	address = *_struct.NewAddress(os.Args[1], os.Args[2])

	netListen, err := net.Listen("tcp", address.ToString())
	if err != nil {
		log.Fatal(err.Error())
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor))

	ExecuteService := service.KeyValueServiceImpl{}
	pb.RegisterKeyValueServiceServer(server, &ExecuteService)

	log.Println("Server Started At ", netListen.Addr())

	if err := server.Serve(netListen); err != nil {
		log.Fatal(err.Error())
	}
}
