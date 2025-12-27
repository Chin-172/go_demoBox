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
	dataList   = make(chan []*pb.DataCenterEntity)
	actionRes  = make(chan int) // 1 - Success ; 0 - Fail
	actionList = make([]string, 0, 100)
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
	// log.Printf("Recv data: username: %s\npassword: %s\n", res.GetUser().Username, res.GetUser().Password)
	if res.GetUser().GetUsername() == "" {
		return nil
	}
	return res
}

func getAccountActionList(identity int, client pb.DataProtocolClient) {
	// var actionList []string
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	var action pb.RPCAction
	action.Action = "GetUserActionList"
	fmt.Printf("identity: %d\n", identity)
	// switch identity {
	// case 0:
	// 	action.Keyword = "0"

	// }

	action.Keyword = strconv.Itoa(identity)

	stream, err := client.ServerStreaming(ctx, &action)
	if err != nil {
		log.Printf("Recv error:%v", err)
		wg.Done()
		return
	}
	// if res.GetAction().GetActionName() != "" {
	// 	actionList = append(actionList, res.GetAction().GetActionName())
	// }
	for {
		resp, err := stream.Recv()
		log.Printf("Recv data: %s\n", resp.GetAction().GetActionName())

		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			// log.Printf("Recv error:%v", err)
			break
		}
		if resp.GetAction().GetActionName() != "" {
			actionList = append(actionList, resp.GetAction().GetActionName())
		}
	}
	wg.Done()
}

func getAllDataList(user *pb.UserEntity, level string, client pb.DataProtocolClient) []*pb.DataCenterEntity {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers
	var data []*pb.DataCenterEntity
	var action pb.RPCAction
	action.Action = "GetAllDataList"
	action.Keyword = level
	action.UserInfo = user
	stream, err := client.ServerStreaming(ctx, &action)
	if err != nil {
		log.Printf("Recv error:%v", err)
		// wg.Done()
		return data
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			log.Printf("Recv error:%v", err)
			break
		}
		// log.Printf("Recv data:%v", resp)

		if resp.GetDataCenter().GetFileName() != "" {
			data = append(data, resp.GetDataCenter())
		}
	}
	return data

}
func actionParse(action string, client pb.DataProtocolClient, user *pb.UserEntity) {
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

	case "Data Center":
		fmt.Println("Data Center In")

		var level string
		fmt.Printf("Please input which level case you want to view: (All/A/B/C):")
		fmt.Scanf("%s\n", &level)
		go func() {
			dataList <- getAllDataList(user, level, client)
		}()

		for i, data := range <-dataList {
			fmt.Printf("%d. %s\n", i+1, data.FileName)
		}
	case "Log":
		fmt.Println("Log In")
	case "Group List":
		fmt.Println("Group List In")
	case "User List":
		fmt.Println("User List In")
	}
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
				actionParse("Login", client, nil)
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
		case 9:
			go func() {
				leave <- true
			}()
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
		var cmd int

		user = <-userInfo
		fmt.Printf("Welcome Back, %s!\n", user.GetUser().Username)
		// GetActionList
		wg.Add(1)
		go func() {
			getAccountActionList(int(user.Id), client)
		}()
		wg.Wait()
		for {

			fmt.Printf("Please input the command \n")
			for i, action := range actionList {
				fmt.Printf("%d. %s\n", i+1, action)
			}
			fmt.Printf("Command: ")

			fmt.Scanf("%d\n", &cmd)
			action := actionList[cmd-1]

			wg.Add(1)
			go func() {
				var userInfo pb.UserEntity
				userInfo = *user.GetUser()
				actionParse(action, client, &userInfo)
			}()
			wg.Wait()
		}

	}
}
