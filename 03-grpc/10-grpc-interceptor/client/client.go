package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"test-google-micro-service/02-grpc/10-grpc-interceptor/proto"
)

func main() {

	intercaptor:=func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
		start:=time.Now()
		err:=invoker(ctx,method,req,reply,cc,opts...)
		fmt.Printf("耗时%s\n",time.Since(start))
		return err
	}
	opt:=grpc.WithUnaryInterceptor(intercaptor)

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), opt)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)

	md := metadata.New(map[string]string{
		"name":    "golang",
		"pasword": "imooc",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "golang"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}
