package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"test-google-micro-service/02-grpc/04-proto-message/proto"
)

func main () {
	conn, err := grpc.Dial("127.0.0.7:8080", grpc.WithInsecure())
	if err != nil {
		panic("grpc.Dial err" + err.Error())
	}

	client := proto.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &proto.HelloRequest{Name: "hello golang", Id: []int32{11, 22}})
	if err != nil {
		panic(err)
	}

	fmt.Println(reply.Message)
	fmt.Println(reply.Data)

	defer conn.Close()
}
