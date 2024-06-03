package connection

import (
	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type GRPCClient struct {
	address  *Address
	conn     *grpc.ClientConn
	Services struct {
		KV   pb.KeyValueServiceClient
		Raft pb.RaftServiceClient
	}
}

func NewClient(_address *Address) (*GRPCClient, error) {
	client := &GRPCClient{address: _address}
	err := client.createConn()
	return client, err
}

func (client *GRPCClient) createConn() (err error) {
	client.conn, err = grpc.NewClient(client.address.ToString(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error Dial %v", err)
	}
	client.Services.KV = pb.NewKeyValueServiceClient(client.conn)
	client.Services.Raft = pb.NewRaftServiceClient(client.conn)
	return err
}

func (client *GRPCClient) SetAddress(_address *Address) error {
	client.conn.Close()
	client.address = _address
	return client.createConn()
}
