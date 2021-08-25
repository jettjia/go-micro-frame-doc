package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func main() {
	// 创建连接
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Printf("连接失败: [%v]\n", err)
		return
	}
	defer conn.Close()
	// 声明客户端
	client := proto_bak.NewGreeterClient(conn)
	rsp, _ := client.SayHello(context.Background(), &proto_bak.HelloRequest{
		Name: "bobby",
		Url:  "https://imooc.com",
		G:    proto_bak.Gender_MALE,
		Mp: map[string]string{
			"name":    "bobby",
			"company": "慕课网",
		},
		AddTime: timestamppb.New(time.Now()),
	})
	fmt.Println(rsp.Message)
}
