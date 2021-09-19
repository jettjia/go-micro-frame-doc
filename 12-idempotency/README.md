# 幂等性

## 介绍

工作当中幂等性是很常见的，我们希望同样的数据提交多次，结果是一样的。

下面介绍一些，需要考虑幂等性的情况。

| 接口类型 | 是否考虑幂等 | 举例                                                         |
| -------- | ------------ | ------------------------------------------------------------ |
| get      | 否           | 多次清，返回结果不一定一样。比如刷新首页，可能广告位的变化等 |
| post     | 是           | 比如前端页面提交了创建订单                                   |
| put      | 不一定       | 不需要考虑的有：修改了商品的价格；需要考虑的有：购物车商品 +1，需要保证幂等性 |
| delete   | 否           | 不用考虑                                                     |



## 幂等性解决方案

1）唯一索引，mysql 增加唯一索引

2） token机制，防止页面重复提交，和前端定义 token 规则；

3） 锁

乐观锁 、悲观锁、分布式锁

4） 提供接口的api保证幂等

比如银联支付接口等

接口必须传两个字段：source(来源)，seq(序列化)，这两个字段在提供方系统里面做唯一索引，防止多次操作。

当第三方调用时，先在本方系统查询一下，是否已经处理过，返回对应的处理结果。



# grpc重试

在client中调用的时候，可以发起多次重试到 grpc中，以此来保证请求；下面有完整代码示例

## proto

```protobuf
syntax = "proto3";
option go_package = ".;proto";
package proto;
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply);
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}

```

```
protoc --go_out=plugins=grpc:./ ./hello.proto
```

## server.go

```go
package main

import (
	"context"
	"fmt"
	"go-micro-frame-doc/12-idempotency/01-grpc/proto"
	"net"
	"time"

	"google.golang.org/grpc"

)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply,
	error) {
	time.Sleep(2 * time.Second)
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func main() {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接收到了一个新的请求")
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

```



## client.go

```go
package main

import (
	"context"
	"fmt"
	"go-micro-frame-doc/12-idempotency/01-grpc/proto"
	"time"

	"google.golang.org/grpc/codes"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"

)

func main() {
	//stream
	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("耗时：%s\n", time.Since(start))
		return err
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3), // 重试3次
		grpc_retry.WithPerRetryTimeout(1 * time.Second), // 每次重试1s
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable), //重试的状态码
	}

	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))
	//这个请求应该多长时间超时， 这个重试应该几次、当服务器返回什么状态码的时候重试
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	conn, err := grpc.Dial("127.0.0.1:50051", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}

```



# 幂等性实践

比如开发中，创建订单需要考虑到幂等性；可以加分布式锁来控制。参考前面的章节