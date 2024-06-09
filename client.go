package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/client"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"

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

	address := NewAddress(os.Args[1], os.Args[2])

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

			// response := &pb.Response{
			// 	RedirectAddress: address.Address,
			// 	Status:          pb.STATUS_REDIRECTED.Enum(),
			// }

			// fmt.Println(response.Status.String() == pb.STATUS_REDIRECTED.String())
			// for response.Status.String() == pb.STATUS_REDIRECTED.String() {
			// 	// fmt.Println("REDIRECT")
			// 	response, err = client.Services.KV.Ping(ctx, &pb.Empty{})
			// 	if err != nil {
			// 		log.Fatalf("Response Error %v", err)
			// 	}
			// }

			// fmt.Printf("%s\n", response.GetValue())

			// function := func() {
			// 	response, err := client.Services.KV.Ping(ctx, &pb.Empty{})
			// 	if err != nil {
			// 		log.Fatalf("Response Error %v", err)
			// 	}

			// 	fmt.Printf("%s\n", response.GetValue())
			// }

			function := func() {
				handler := func(address Address) *pb.Response {
					client.SetAddress(&address)
					response, err := client.Services.KV.Ping(ctx, &pb.Empty{})
					if err != nil {
						log.Fatalf("Response Error %v", err)
					}
					return response
				}
				response := RedirectHanlder(address, handler)
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
				fmt.Printf("\"%s\"\n", response.GetValue())
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
				fmt.Printf("\"%s\"\n", response.GetValue())
			}

			if enableTime {
				TimeWrap(function)
			} else {
				function()
			}

		case "strln":
			if len(command) != 2 {
				fmt.Println("Invalid put command. Format: strln <key>")
				continue
			}

			ctx := context.Background()

			function := func() {
				response, err := client.Services.KV.StrLn(ctx, &pb.KeyRequest{Key: command[1]})
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

		case "del":
			if len(command) != 2 {
				fmt.Println("Invalid put command. Format: del <key>")
				continue
			}

			ctx := context.Background()

			function := func() {
				response, err := client.Services.KV.Del(ctx, &pb.KeyRequest{Key: command[1]})
				if err != nil {
					log.Fatalf("Response Error %v", err)
				}
				fmt.Printf("\"%s\"\n", response.GetValue())
			}

			if enableTime {
				TimeWrap(function)
			} else {
				function()
			}

		case "append":
			if len(command) != 3 {
				fmt.Println("Invalid put command. Format: append <key> <value>")
				continue
			}

			ctx := context.Background()

			function := func() {
				response, err := client.Services.KV.Append(ctx, &pb.KeyValueRequest{Key: command[1], Value: command[2]})
				if err != nil {
					log.Fatalf("Response Error %v", err)
				}
				fmt.Printf("\"%s\"\n", response.GetValue())
			}

			if enableTime {
				TimeWrap(function)
			} else {
				function()
			}

		case "request_log":
			if len(command) != 2 {
				fmt.Println("Invalid put command. Format: request_log")
				continue
			}

			ctx := context.Background()

			function := func() {
				// TODO: Change To RequestLog
				response, err := client.Services.KV.RequestLog(ctx, &pb.Empty{})
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
			fmt.Println("Error: Unknown command")
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

func RedirectHanlder(address *Address, function func(Address) *pb.Response) *pb.Response {
	response := &pb.Response{
		RedirectAddress: address.Address,
		Status:          pb.STATUS_REDIRECTED.Enum(),
	}

	for response.Status.String() == pb.STATUS_REDIRECTED.String() {
		fmt.Println("REDIRECT")
		response = function(Address{Address: response.RedirectAddress})
	}

	return response
}
