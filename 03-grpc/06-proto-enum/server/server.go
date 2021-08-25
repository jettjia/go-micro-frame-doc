package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"test-google-micro-service/02-grpc/06-proto-enum/proto"
)

type Server struct {
}

func (s *Server) SayHello(ct context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	rs := []*proto.Result{
		&proto.Result{Code: "11", Msg: "success"},
		&proto.Result{Code: "00", Msg: "fail"},
	}

	fmt.Println(fmt.Sprintf("Gender is %v", req.Sex))

	return &proto.HelloReply{
		Message: "hello " + req.Name,
		Data:    rs,
	}, nil
}

func main() {
	g := grpc.NewServer()
	s := Server{}
	proto.RegisterGreeterServer(g, &s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":8080"))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}
	g.Serve(lis)
}
