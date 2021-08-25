package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"

	"google.golang.org/grpc/metadata"

	"test-google-micro-service/02-grpc/10-grpc-interceptor/proto"
)

type Server struct {
}

func (this *Server) SayHello(ctx context.Context, res *proto.HelloRequest) (resp *proto.HelloReply, err error) {
	// 接收metadata信息
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md, ok)
	return &proto.HelloReply{
		Message: fmt.Sprintf("name is %s", res.Name),
	}, nil
}

func main() {

	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接受到一个请求")
		res, err := handler(ctx, req)
		fmt.Println("请求已经完成")
		return res, err
	}

	opt := grpc.UnaryInterceptor(interceptor)
	g := grpc.NewServer(opt)
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
