package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"test-google-micro-service/02-grpc/03-proto-import/proto"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		panic("net.Dial err " + err.Error())
	}

	client := proto.NewGreeterClient(conn)

	reply, err := client.SayHello(context.Background(), &proto.HelloRequest{Name: "hello name", Id: []int32{1, 2, 3}})
	if err != nil {
		panic(err)
	}
	fmt.Println(reply.Message)

	pong, err := client.Ping(context.Background(), &proto.Empty{Name: "empty name"})
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)

	defer conn.Close()
}
