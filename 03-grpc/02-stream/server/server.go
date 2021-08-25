package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"test-google-micro-service/02-grpc/02-stream/proto"
)

const PORT = ":50052"

type server struct {}

// 服务器流
func (s *server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++

		_ = res.Send(&proto.StreamResData{
			Data: fmt.Sprintf("%v", time.Now().Unix()),
		})
		time.Sleep(time.Second)

		if i > 10 {
			break
		}
	}

	return nil
}

func (s *server) PutStream(cliStr proto.Greeter_PutStreamServer) error {
	for {
		if data, err := cliStr.Recv(); err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(data.Data)
		}
	}

	return nil
}

func (s *server) AllStream(allStr proto.Greeter_AllStreamServer) error{
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			data, _ := allStr.Recv()
			fmt.Println("收到客户端的消息:" + data.Data)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			_ = allStr.Send(&proto.StreamResData{
				Data: "我是服务器",
			})
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	listener, _ := net.Listen("tcp", PORT)
	s := grpc.NewServer()

	proto.RegisterGreeterServer(s, &server{})
	err := s.Serve(listener)

	if err != nil {
		panic(err)
	}
}
