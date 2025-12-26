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
	port       = flag.Int("Port", 50002, "Server Port")
	login      = make(chan int)
	userInfo   = make(chan *pb.DataEntity)
	actionRes  = make(chan int) // 1 - Success ; 0 - Fail
	actionList = make(chan []string)
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

func callGetUserInfo(username string, client pb.DataProtocolClient) *pb.DataEntity {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	var action pb.RPCAction
	action.Action = "GetUserInfo"
	action.Keyword = username

	res, err := client.ServerSending(ctx, &action)
	if err != nil {
		log.Printf("Recv error:%v", err)
		return nil
	}
	log.Printf("Recv data: username: %s\npassword: %s\n", res.GetUser().Username, res.GetUser().Password)
	if res.GetUser().GetUsername() == "" {
		return nil
	}
	// wg.Done()
	return res
}

func actionParse(action string, client pb.DataProtocolClient) {
	switch action {
	case "Login":
		var username string
		var pwd string
		fmt.Printf("Username: ")
		fmt.Scanln(&username)

		fmt.Printf("Password: ")
		fmt.Scanln(&pwd)
		for i := 1; i < 3; i++ {
			var user *pb.DataEntity
			// wg.Add(1)
			go func() {
				userInfo <- callGetUserInfo(username, client)
			}()
			// wg.Wait()

			user = <-userInfo
			if user != nil {
				if user.GetUser().GetPassword() != pwd {
					fmt.Printf("Password Incorrect! Please retry again\n")
					fmt.Printf("Password: ")
					fmt.Scanln(&pwd)

					if i == 2 {
						go func() {
							actionRes <- 0
						}()
					}
				} else {
					fmt.Printf("Account Login Success!\n")
					// cmd = 9
					go func() {
						login <- 1
					}()
					go func() {
						userInfo <- user
					}()
					go func() {
						actionRes <- 1
					}()
					break
				}
			} else {
				fmt.Printf("Username not found!\n")
				go func() {
					actionRes <- 0
				}()
				break
			}
		}
	}
	fmt.Println("End")
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
		fmt.Printf("Please input the command (1-System Login; 2-callGetUserInfo; 9-Exit): ")
		fmt.Scanf("%d\n", &cmd)
		switch cmd {
		case 1:
			wg.Add(1)
			go func() {
				actionParse("Login", client)
			}()
			wg.Wait()
			res := <-actionRes
			fmt.Printf("%d\n", res)
			var release sync.WaitGroup

			if res == 0 {
				release.Go(func() {
					fmt.Printf("1. cmd:%d\n", cmd)
					isLeave := cmd == 9
					leave <- isLeave
				})
			} else {
				release.Go(func() {
					leave <- true
				})
			}
		case 2:
		default:
			go func() {
				login <- 0
			}()
		}
		fmt.Printf("cmd:%d\n", cmd)

		if <-leave {
			break
		}
	}

	goLoginSystem := <-login
	fmt.Println(goLoginSystem)
	if 1 == goLoginSystem {
		var user *pb.DataEntity
		user = <-userInfo
		fmt.Printf("Welcome Back, %s!\n", user.GetUser().Username)
		for {
			var cmd int
			fmt.Printf("Please input the command (1-System Login; 2-callGetUserInfo; 9-Exit): ")
			fmt.Scanf("%d\n", &cmd)
			switch cmd {
			case 1:
				wg.Add(1)
				go func() {
					actionParse("Login", client)
				}()
				wg.Wait()
				res := <-actionRes
				fmt.Printf("%d\n", res)
				var release sync.WaitGroup

				if res == 0 {
					release.Go(func() {
						fmt.Printf("1. cmd:%d\n", cmd)
						isLeave := cmd == 9
						leave <- isLeave
					})
				} else {
					release.Go(func() {
						leave <- true
					})
				}
			case 2:
			default:
				go func() {
					login <- 0
				}()
			}

		}
	}
}
