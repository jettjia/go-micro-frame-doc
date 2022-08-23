# 官网
etcd.io

文档: https://etcd.io/docs/v3.5/quickstart/
https://doczhcn.gitbook.io/etcd/index/index/interacting_v3

# 案例模块

```shell
01-kv: 键值对的 put/get/delete
02-lease: 租约
03-watcher: watch机制
04-op: etcd的op操作
05-txn: 事务
06-lock: etcd的分布式锁
  参:
  https://github.com/etcd-io/etcd/blob/v3.3.10/clientv3/concurrency/example_mutex_test.go
  https://pandaychen.github.io/2019/10/24/ETCD-DISTRIBUTED-LOCK/
  https://chai2010.cn/advanced-go-programming-book/ch6-cloud/ch6-02-lock.html
07-election: 选leader
  https://github.com/etcd-io/etcd/blob/v3.3.10/clientv3/concurrency/example_election_test.go
```