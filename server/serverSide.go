// 模擬伺服器推送資料到Client 端上
// 輸入數字 1 作定時執行的模擬

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/chin-172/go_demoBox/proto" // go.mod 的module name + 目錄
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	dbData "github.com/chin-172/go_demoBox/db"
)

type data struct {
	pb.UnimplementedDataProtocolServer
}

var (
	port     = flag.Int("Port", 50002, "Listening port")
	userRecv = make(chan []dbData.User)
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
				Auth:     uint32(user.Auth),
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

		userEntity := &pb.UserEntity{
			Id:       uint32(userInfo.ID),
			Username: userInfo.Username,
			Password: userInfo.Password,
			Auth:     uint32(userInfo.Auth),
		}

		passingData.Id = uint32(1)
		passingData.Data = &pb.DataEntity_User{
			User: userEntity,
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
