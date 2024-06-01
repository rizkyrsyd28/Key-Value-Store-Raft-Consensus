package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	_struct "github.com/Sister20/if3230-tubes-dark-syster/lib/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
	address   _struct.Address
)

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
		if os.Args[1] == "time" {
			enableTime = true
		}
	}

	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Failed to convert address port: %v", err)
	}

	address = *_struct.NewAddress(os.Args[1], port)

	conn, err := grpc.NewClient(address.IP+":"+strconv.Itoa(address.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error Dial %v", err)
	}
	defer conn.Close()

	client := pb.NewExecuteServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		// fmt.Print("Enter command: \n> ")
		fmt.Printf("%s:%d> ", address.IP, address.Port)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		// Remove newline character
		input = strings.TrimSpace(input)

		command := strings.Fields(input)

		switch command[0] {
		case "ping":
			if len(command) != 1 {
				fmt.Println("Invalid put command. Format: ping")
				continue
			}

			ctx := context.Background()

			if enableTime {
				startTime = time.Now()
			}

			response, err := client.Ping(ctx, &pb.Empty{})
			if err != nil {
				log.Fatalf("Response Error %v", err)
			}

			if enableTime {
				endTime = time.Now()
				duration = endTime.Sub(startTime)
			}

			fmt.Printf("%s\n", response.GetValue())

			if enableTime {
				fmt.Printf("Response Time: %v\n", duration)
			}

		case "get":
			if len(command) != 2 {
				fmt.Println("Invalid put command. Format: get <key>")
				continue
			}

			ctx := context.Background()

			if enableTime {
				startTime = time.Now()
			}

			response, err := client.Get(ctx, &pb.KeyRequest{Key: command[1]})
			if err != nil {
				log.Fatalf("Response Error %v", err)
			}

			if enableTime {
				endTime = time.Now()
				duration = endTime.Sub(startTime)
			}

			fmt.Printf("%s\n", response.GetValue())

			if enableTime {
				fmt.Printf("Response Time: %v\n", duration)
			}

		case "put":
			if len(command) != 3 {
				fmt.Println("Invalid put command. Format: put <key> <value>")
				continue
			}

			ctx := context.Background()

			if enableTime {
				startTime = time.Now()
			}

			response, err := client.Get(ctx, &pb.KeyRequest{Key: command[1]})
			if err != nil {
				log.Fatalf("Response Error %v", err)
			}

			if enableTime {
				endTime = time.Now()
				duration = endTime.Sub(startTime)
			}

			fmt.Printf("%s\n", response.GetValue())

			if enableTime {
				fmt.Printf("Response Time: %v\n", duration)
			}

		case "quit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}