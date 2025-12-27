// 模擬伺服器推送資料到Client 端上
// 輸入數字 1 作定時執行的模擬

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	pb "github.com/chin-172/go_demoBox/proto" // go.mod 的module name + 目錄
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	dbData "github.com/chin-172/go_demoBox/db"
)

type data struct {
	pb.UnimplementedDataProtocolServer
}

var (
	port       = flag.Int("Port", 50002, "Listening port")
	userRecv   = make(chan []dbData.User)
	actionList = make(chan []string)
	dataList   = make(chan []dbData.DataCenter)
	authPass   = make(chan bool)
)

func (d *data) ServerStreaming(request *pb.RPCAction, stream pb.DataProtocol_ServerStreamingServer) error {
	// ctx := stream.Context() // 從stream 中獲取context
	fmt.Printf("%x online\n", stream.Context().Done())

	action := request.GetAction()
	// Get data from DB
	dbData.DBconn()
	switch action {
	case "GetAllUser":
		go func() {
			userRecv <- dbData.GetAllUser()
		}()

		// Send data throug stream
		for i, user := range <-userRecv {
			userEntity := &pb.UserEntity{
				Id:       uint32(user.ID),
				Username: user.Username,
				Password: user.Password,
				Identity: uint32(user.Identity),
				GroupID:  uint32(user.GroupID),
			}

			var passingData pb.DataEntity
			passingData.Id = uint32(i)
			passingData.Data = &pb.DataEntity_User{
				User: userEntity,
			}

			err := stream.Send(&passingData)

			if err != nil {
				fmt.Printf("Send error:%v\n", err)
				return err
			}
		}

	case "GetUserActionList":
		auth, err := strconv.Atoi(request.GetKeyword())
		if err != nil {
			fmt.Printf("Int Parse error:%v\n", err)
			return err
		}
		go func() {
			actionList <- dbData.GetUserActionList(auth)
		}()

		for i, action := range <-actionList {
			fmt.Printf("%d %s\n", i, action)
			actionEntity := &pb.ActionEntity{
				Id:         uint32(i),
				ActionName: action,
			}
			var passingData pb.DataEntity
			passingData.Id = uint32(i)
			passingData.Data = &pb.DataEntity_Action{
				Action: actionEntity,
			}

			err := stream.Send(&passingData)

			if err != nil {
				fmt.Printf("Send error:%v\n", err)
				return err
			}
		}

	case "GetAllDataList":
		// Check authorization of data query
		var caseLevel []string
		if strings.Contains(request.GetKeyword(), "+") {
			caseLevel = strings.Split(request.GetKeyword(), "+")
		} else {
			caseLevel = append(caseLevel, request.GetKeyword())
		}
		fmt.Printf("caseLevel: %s\n", caseLevel)

		go func() {
			authPass <- dbData.CheckAuth(int(request.GetUserInfo().GroupID), caseLevel)
		}()

		pass := <-authPass
		if !pass {
			dbData.CloseDBConn()
			return fmt.Errorf("User have no authority to check this case level data")
		}
		fmt.Printf("Auth Pass!\n")
		go func() {
			dataList <- dbData.GetAllDataList(caseLevel)
		}()

		for i, data := range <-dataList {
			dataEntity := &pb.DataCenterEntity{
				Id:        uint32(i),
				FileName:  data.FileName,
				Year:      uint32(data.Year),
				Signatory: data.Signatory,
				FileType:  data.FileType,
				AddDt:     data.AddDt.Unix(),
				Operator:  data.Operator,
			}
			var passingData pb.DataEntity
			passingData.Id = uint32(i)
			passingData.Data = &pb.DataEntity_DataCenter{
				DataCenter: dataEntity,
			}

			err := stream.Send(&passingData)

			if err != nil {
				fmt.Printf("Send error:%v\n", err)
				return err
			}
		}
	}
	fmt.Printf("%x close\n", stream.Context().Done())
	dbData.CloseDBConn()

	return nil
}

// 使用控制流ctx 控制超時等操作
func (d *data) ServerSending(ctx context.Context, request *pb.RPCAction) (*pb.DataEntity, error) {
	action := request.GetAction()

	dbData.DBconn()
	var passingData pb.DataEntity

	switch action {
	case "GetUserInfo":
		// var userEntity pb.DataEntity_User

		userInfo := dbData.GetUserInfo(request.GetKeyword())
		if userInfo.Username == "" {
			passingData.Id = uint32(1)
			passingData.Data = &pb.DataEntity_User{
				User: nil,
			}
		} else {
			userEntity := &pb.UserEntity{
				Id:       uint32(userInfo.ID),
				Username: userInfo.Username,
				Password: userInfo.Password,
				Identity: uint32(userInfo.Identity),
			}

			passingData.Id = uint32(1)
			passingData.Data = &pb.DataEntity_User{
				User: userEntity,
			}
		}

	default:
		passingData.Id = uint32(1)
		passingData.Data = nil
	}
	fmt.Printf("%x close\n", ctx.Done())
	dbData.CloseDBConn()
	return &passingData, nil
}
func establishConn() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server) // allows to use grpcui
	pb.RegisterDataProtocolServer(server, &data{})
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
func main() {
	for {
		var cmd int
		fmt.Printf("Please input the command (1-Open Port; 9-Exit): ")

		fmt.Scanf("%d\n", &cmd)

		if cmd == 9 {
			break
		} else if cmd == 1 {
			// run push
			// var port int
			// fmt.Printf("Please input the listening port (localhost): ")
			// fmt.Scanf("%d\n", &port)
			go func() {
				establishConn()
			}()
		}
	}

	fmt.Println("Excution done")
}
