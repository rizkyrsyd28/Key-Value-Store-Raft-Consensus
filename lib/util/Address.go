package util

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
)

type Address struct {
	*pb.Address
}

func NewAddress(ip string, port string) *Address {
	_port, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		log.Fatalf("Failed to convert address port: %v", err)
	}
	return &Address{
		&pb.Address{
			IP:   ip,
			Port: uint32(_port),
		},
	}
}

func (address *Address) ToString() string {
	return fmt.Sprintf("%s:%d", address.IP, address.Port)
}

func (address *Address) Iterator() <-chan interface{} {
	iter := make(chan interface{}, 2)
	iter <- address.IP
	iter <- address.Port
	close(iter)
	return iter
}

func (address *Address) IsEqual(other *Address) bool {
	return address.IP == other.IP && address.Port == other.Port
}

func (address *Address) IsNotEqual(other *Address) bool {
	return address.IP != other.IP || address.Port != other.Port
}

// func ToAddress(pbAddress pb.Address) *Address {
// 	return &Address{
// 		Address: pbAddress,
// 	}
// }
