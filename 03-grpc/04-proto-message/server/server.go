package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"test-google-micro-service/02-grpc/04-proto-message/proto"
)

type Server struct {
}

func (s *Server) SayHello(ct context.Context, req *proto.HelloRequest) (resp *proto.HelloReply, err error) {
	rs := []*proto.HelloReply_Result{
		&proto.HelloReply_Result{Code: "11", Msg: "success"},
		&proto.HelloReply_Result{Code: "22", Msg: "fail"},
	}
	return &proto.HelloReply{
		Message: fmt.Sprintf("name is %s", req.Name),
		Data:    rs,
	}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic("net.Listen err " + err.Error())
	}

	_ = g.Serve(listener)
}
