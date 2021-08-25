package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"test-google-micro-service/02-grpc/01-helloword/proto"
)

func main() {
	conn, _ := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())

	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	rsp, _ :=c.SayHello(context.Background(),
		&proto.HelloRequest{Name:"golang"},
	)

	fmt.Println(rsp.Message)
}
