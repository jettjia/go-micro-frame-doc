package main

import (
	"context"
	"google.golang.org/grpc"
	"net"

	"test-google-micro-service/02-grpc/01-helloword/proto"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: "hello, " + req.Name,
	}, nil
}

func main() {
	g := grpc.NewServer()

	proto.RegisterGreeterServer(g, &Server{})

	Listener, _ := net.Listen("tcp", "0.0.0.0:8080")

	_ = g.Serve(Listener)
}
