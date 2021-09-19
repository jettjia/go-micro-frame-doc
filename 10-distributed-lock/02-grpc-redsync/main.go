package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"go-micro-frame-doc/10-distributed-lock/02-grpc-redsync/global"
	"go-micro-frame-doc/10-distributed-lock/02-grpc-redsync/handler"
	"go-micro-frame-doc/10-distributed-lock/02-grpc-redsync/initialize"
	"go-micro-frame-doc/10-distributed-lock/02-grpc-redsync/proto"
	"go-micro-frame-doc/10-distributed-lock/02-grpc-redsync/utils"
	"go-micro-frame-doc/10-distributed-lock/02-grpc-redsync/utils/register/consul"
)

func main() {
	// 判断是否生成 随机的 微服务端口号
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口")
	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0{
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	// 初始化 logger
	initialize.InitLogger()

	//初始化配置文件
	initialize.InitConfig()

	// 初始化db
	initialize.InitDB()

	// 初始化 redsyncLock
	initialize.InitRedsyncLock()

	zap.S().Info(global.ServerConfig)

	/////////////////////////////////
	// 启动grpc，并注册到 consul
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口： %d", *Port)

	/////////////////////////////////

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功")
	}

}
