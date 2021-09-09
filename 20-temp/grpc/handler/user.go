package handler

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-micro-module/20-temp/grpc/global"
	"go-micro-module/20-temp/grpc/model"
	"go-micro-module/20-temp/grpc/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToRsponse(user model.User) proto.UserInfoResponse {
	//在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	//这里要搞清， 哪些字段是有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// 获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	fmt.Println("我被调用了")
	fmt.Println("我调用的数据库配置是", global.ServerConfig.MysqlInfo)

	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// 这里模拟，有调用第三方的服务
	span, ctx := opentracing.StartSpanFromContext(ctx,"main")
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	span.Finish()

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}

// 这里本来是根据用户id 获取用户信息的
// 在这里增加一个 分布式的 redis 锁，来测试；注意这里只是模拟，实际业务一般是可以在控制库存 加减的时候
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	// 设置锁
	mutex := global.RedsyncLock.NewMutex(fmt.Sprintf("userid_", req.Id))
	if err := mutex.Lock(); err != nil {
		return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
	}

	result := global.DB.First(&user, req.Id)
	// 释放分布式锁
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
	}

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}
