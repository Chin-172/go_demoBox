// Simulate the client trigger gRPC from server to get the data
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	pb "github.com/chin-172/go_demoBox/proto" // go.mod 的module name + 目錄

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("Port", 50002, "Server Port")
)

var wg sync.WaitGroup

func callGetAllUser(client pb.DataProtocolClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	// emptyReq := &emptypb.Empty{}

	var action pb.RPCAction
	action.Action = "GetAllUser"

	stream, err := client.ServerStreaming(ctx, &action)
	if err != nil {
		log.Fatalf("ServerStreaming stream err: %v", err)
		return
	}

	// Receive the streaming from server
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			// log.Printf("Recv error:%v", err)
			break
		}
		log.Printf("Recv data: %s\n", resp.GetUser().Username)
	}
	wg.Done()

}

func callGetUserInfo(client pb.DataProtocolClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	var action pb.RPCAction
	action.Action = "GetUserInfo"
	action.Keyword = "Moon"

	res, err := client.ServerSending(ctx, &action)
	if err != nil {
		log.Printf("Recv error:%v", err)
	}
	log.Printf("Recv data:\nusername: %s\npassword: %s\n", res.GetUser().Username, res.GetUser().Password)

	wg.Done()

}
func main() {
	flag.Parse()
	connStr := "localhost:" + strconv.Itoa(*port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDataProtocolClient(conn)
	leave := make(chan bool)
	for {
		var cmd int
		fmt.Printf("Please input the command (1-callGetAllUser; 2-callGetUserInfo; 9-Exit): ")
		fmt.Scanf("%d\n", &cmd)
		wg.Add(1)
		switch cmd {
		case 1:
			go func() {
				callGetAllUser(client)
			}()
			wg.Wait()

		case 2:
			go func() {
				callGetUserInfo(client)
			}()
			wg.Wait()
		default:
			wg.Done()
		}

		go func() {
			leave <- cmd == 9
		}()
		if <-leave {
			break
		}
	}

}
