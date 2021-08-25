package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"test-google-micro-service/02-grpc/13-grpc-error/proto"
)

func main() {
	//stream
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	//go语言推荐的是返回一个error和一个正常的信息
	_, err = c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a status error
			panic("解析error失败")
		}
		fmt.Println(st.Message())
		fmt.Println(st.Code())
	}
	//fmt.Println(r.Message)
}
