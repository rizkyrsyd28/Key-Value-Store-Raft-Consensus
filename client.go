package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/client"
	_struct "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
)

//var address _struct.Address

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run client.go <server ip> <server port> <time?>")
		return
	} else {
		if os.Args[2] == "time" {
			fmt.Println("Usage: go run client.go <server ip> <server port> <time?>")
			return
		}
	}

	enableTime := false

	if len(os.Args) == 4 {
		if os.Args[3] == "time" {
			enableTime = true
		}
	}

	address := _struct.NewAddress(os.Args[1], os.Args[2])

	client, err := client.NewClient(address)
	if err != nil {
		log.Fatalf("Error Dial %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s:%d> ", address.IP, address.Port)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		input = strings.TrimSpace(input)

		command := strings.Fields(input)

		switch command[0] {
		case "ping":
			if len(command) != 1 {
				fmt.Println("Invalid put command. Format: ping")
				continue
			}

			ctx := context.Background()

			function := func() {
				response, err := client.Services.KV.Ping(ctx, &pb.Empty{})
				if err != nil {
					log.Fatalf("Response Error %v", err)
				}

				fmt.Printf("%s\n", response.GetValue())
			}

			if enableTime {
				TimeWrap(function)
			} else {
				function()
			}

		case "get":
			if len(command) != 2 {
				fmt.Println("Invalid put command. Format: get <key>")
				continue
			}

			ctx := context.Background()

			function := func() {
				response, err := client.Services.KV.Get(ctx, &pb.KeyRequest{Key: command[1]})
				if err != nil {
					log.Fatalf("Response Error %v", err)
				}
				fmt.Printf("%s\n", response.GetValue())
			}

			if enableTime {
				TimeWrap(function)
			} else {
				function()
			}

		case "set":
			if len(command) != 3 {
				fmt.Println("Invalid put command. Format: put <key> <value>")
				continue
			}

			ctx := context.Background()

			function := func() {
				response, err := client.Services.KV.Set(ctx, &pb.KeyValueRequest{Key: command[1], Value: command[2]})
				if err != nil {
					log.Fatalf("Response Error %v", err)
				}
				fmt.Printf("%s\n", response.GetValue())
			}

			if enableTime {
				TimeWrap(function)
			} else {
				function()
			}

		case "quit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}

func TimeWrap(function func()) {

	startTime := time.Now()

	function()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Response Time: %v\n", duration)
}
