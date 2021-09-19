package handler

import (
	"context"
	"fmt"

	"go-micro-frame-doc/08-nacos/02-grpc-nacos/global"
	"go-micro-frame-doc/08-nacos/02-grpc-nacos/model"
	"go-micro-frame-doc/08-nacos/02-grpc-nacos/proto"
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

	fmt.Println("输出，db连接是否变了", global.ServerConfig.MysqlInfo)

	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}
