package main

import (
	"fmt"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/server"
	_struct "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run server.go <server ip> <server port>")
		return
	}

	address := *_struct.NewAddress(os.Args[1], os.Args[2])

	serverInstance := server.NewServer(&address)

	serverInstance.Serve()
}
