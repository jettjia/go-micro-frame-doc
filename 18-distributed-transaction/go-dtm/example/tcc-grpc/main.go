package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yedf/dtmcli"
	"github.com/yedf/dtmcli/dtmimp"
	"github.com/yedf/dtmgrpc"
	"github.com/yedf/dtmgrpc-go-sample/busi"
	"github.com/yedf/dtmgrpc/dtmgimp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// todo 调试未成功

// DtmGrpcServer dtm grpc service address
const DtmGrpcServer = "localhost:58080"

// BusiGrpcPort 1
const BusiGrpcPort = 50589

// BusiGrpc busi service grpc address
var BusiGrpc string = fmt.Sprintf("localhost:%d", BusiGrpcPort)

func handleGrpcBusiness(in *busi.BusiReq, result1 string, busi string) error {
	res := dtmimp.OrString(result1, dtmcli.ResultSuccess)
	dtmimp.Logf("grpc busi %s result: %s", busi, res)
	if res == dtmcli.ResultSuccess {
		return nil
	} else if res == dtmcli.ResultFailure {
		return status.New(codes.Aborted, "FAILURE").Err()
	}
	return status.New(codes.Internal, fmt.Sprintf("unknow result %s", res)).Err()
}

// busiServer is used to implement helloworld.GreeterServer.
type busiServer struct {
	busi.UnimplementedBusiServer
}

// GrpcStartup for grpc
func GrpcStartup() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", BusiGrpcPort))
	dtmimp.FatalIfError(err)
	s := grpc.NewServer(grpc.UnaryInterceptor(dtmgimp.GrpcServerLog))
	s.RegisterService(&busi.Busi_ServiceDesc, &busiServer{})
	go func() {
		dtmimp.Logf("busi grpc listening at %v", lis.Addr())
		err := s.Serve(lis)
		dtmimp.FatalIfError(err)
	}()
	time.Sleep(100 * time.Millisecond)
}

func (s *busiServer) TransInRevert(ctx context.Context, in *busi.BusiReq) (*busi.BusiReply, error) {
	return &busi.BusiReply{}, handleGrpcBusiness(in, "", dtmimp.GetFuncName())
}

func (s *busiServer) TransOutRevert(ctx context.Context, in *busi.BusiReq) (*busi.BusiReply, error) {
	return &busi.BusiReply{}, handleGrpcBusiness(in, "", dtmimp.GetFuncName())
}

func (s *busiServer) TransInConfirm(ctx context.Context, in *busi.BusiReq) (*busi.BusiReply, error) {
	return &busi.BusiReply{}, handleGrpcBusiness(in, "", dtmimp.GetFuncName())
}

func (s *busiServer) TransOutConfirm(ctx context.Context, in *busi.BusiReq) (*busi.BusiReply, error) {
	return &busi.BusiReply{}, handleGrpcBusiness(in, "", dtmimp.GetFuncName())
}

func (s *busiServer) TransInTcc(ctx context.Context, in *busi.BusiReq) (*busi.BusiReply, error) {
	return &busi.BusiReply{}, handleGrpcBusiness(in, in.TransInResult, dtmimp.GetFuncName())
}

func (s *busiServer) TransOutTcc(ctx context.Context, in *busi.BusiReq) (*busi.BusiReply, error) {
	return &busi.BusiReply{}, handleGrpcBusiness(in, in.TransOutResult, dtmimp.GetFuncName())
}

func main() {
	GrpcStartup()
	dtmimp.Logf("tcc simple transaction begin")
	gid := dtmgrpc.MustGenGid(DtmGrpcServer)
	err := dtmgrpc.TccGlobalTransaction(DtmGrpcServer, gid, func(tcc *dtmgrpc.TccGrpc) error {
		req := &busi.BusiReq{Amount: 30}
		reply := busi.BusiReply{}
		err := tcc.CallBranch(req, BusiGrpc+"/busi.Busi/TransOutTcc", BusiGrpc+"/busi.Busi/TransOutConfirm", BusiGrpc+"/busi.Busi/TransOutRevert", &reply)
		if err != nil {
			return err
		}
		err = tcc.CallBranch(req, BusiGrpc+"/busi.Busi/TransInTcc", BusiGrpc+"/busi.Busi/TransInConfirm", BusiGrpc+"/busi.Busi/TransInRevert", &reply)
		return err
	})
	dtmimp.FatalIfError(err)
	time.Sleep(20 * time.Second)
}
