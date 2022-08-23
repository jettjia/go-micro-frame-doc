package _1_kv

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func Get() {
	var (
		config  clientv3.Config
		err     error
		client  *clientv3.Client
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"10.4.7.71:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}

	// 建立一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs, getResp.Count)
	}
}

func GetPrefix() {
	var (
		config  clientv3.Config
		err     error
		client  *clientv3.Client
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"10.4.7.71:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}

	// 建立一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	// 写入第二个键值对
	kv.Put(context.TODO(), "/cron/jobs/job2", "job2-content")

	// 读取/cron/jobs/为前缀的所有key
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else { // 获取成功, 我们遍历所有的kvs
		fmt.Println(getResp.Kvs, getResp.Count)
	}
}
