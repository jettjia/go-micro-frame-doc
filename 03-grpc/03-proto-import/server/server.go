package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"

	"test-google-micro-service/02-grpc/03-proto-import/proto"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: fmt.Sprintf("name is %s id is %v", req.Name, req.Id),
	}, nil
}

func (s *Server) Ping(ctx context.Context, res *proto.Empty) (resp *proto.Pong, err error) {
	return &proto.Pong{
		Name: fmt.Sprintf("pong name is %s", res.Name),
	}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	listener, err := net.Listen("tcp", fmt.Sprintf(":8080"))
	if err != nil {
		panic("failed to listen" + err.Error())
	}
	_ = g.Serve(listener)
}
