package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"test-google-micro-service/02-grpc/05-proto-message-out/proto"
)

type Server struct {

}

func (this *Server) SayHello(ctx context.Context,res *proto.HelloRequest)(resp *proto.HelloReply,err error){
	rs:=[]*proto.Result{
		&proto.Result{Code: "11",Msg: "success"},
		&proto.Result{Code: "00",Msg: "fail"},
	}
	return &proto.HelloReply{
		Message: fmt.Sprintf("name is %s",res.Name),
		Data: rs,
	},nil
}

func main()  {
	g:=grpc.NewServer()
	s := Server{}
	proto.RegisterGreeterServer(g,&s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":8080"))
	if err != nil {
		panic("failed to listen: "+err.Error())
	}

	g.Serve(lis)
}