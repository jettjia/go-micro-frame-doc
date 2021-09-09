package initialize

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"

	"go-micro-module/17-es/02-grpc/global"
	"go-micro-module/17-es/02-grpc/model"
)

func InitEs() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "mxshop", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(
		elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetBasicAuth(global.ServerConfig.EsInfo.User, global.ServerConfig.EsInfo.Password),
		elastic.SetTraceLog(logger),
	)
	if err != nil {
		panic(err)
	}

	//新建mapping和index
	exists, err := global.EsClient.IndexExists(model.EsUser{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		_, err = global.EsClient.CreateIndex(model.EsUser{}.GetIndexName()).BodyString(model.EsUser{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
