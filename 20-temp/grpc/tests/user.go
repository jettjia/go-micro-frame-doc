package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"go-micro-module/20-temp/grpc/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init(){
	var err error
	conn, err = grpc.Dial("10.4.7.71:45995", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList(){
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		if err != nil {
			panic(err)
		}
	}
}

func TestGetUserById()  {
	rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}


func main() {
	Init()
	TestGetUserList()
	//TestGetUserById()

	conn.Close()
}