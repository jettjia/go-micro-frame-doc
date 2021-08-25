package main
import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"test-google-micro-service/02-grpc/05-proto-message-out/proto"
)
func main()  {
	conn,err:=grpc.Dial("127.0.0.1:8080",grpc.WithInsecure())
	if err!=nil {
		fmt.Println("net.Dial err"+err.Error())
	}
	c := proto.NewGreeterClient(conn)

	r,err := c.SayHello(context.Background(),&proto.HelloRequest{Name:"hello name",Id: []int32{1,2,3}})

	if err!=nil{
		panic(err)
	}
	fmt.Println(r.Message)
	fmt.Println(r.Data)

	defer conn.Close()
}