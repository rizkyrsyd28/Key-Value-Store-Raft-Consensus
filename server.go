package main

import (
	"fmt"
	"os"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/server"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

func main() {
	var isContact bool = false
	var address, contactAddress Address

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run server.go <server ip> <server port> <?contact ip> <?contact port>")
		return
	}

	address = *NewAddress(os.Args[1], os.Args[2])

	if len(os.Args) == 5 {
		contactAddress = *NewAddress(os.Args[3], os.Args[4])
		isContact = true
		fmt.Println("Node Running With Contact Address")
	}

	serverInstance := server.NewServer(&address, isContact, &contactAddress)

	serverInstance.Serve()
}
