# 安装go

```shell
wget https://studygolang.com/dl/golang/go1.15.14.linux-amd64.tar.gz
tar -zxvf go1.15.14.linux-amd64.tar.gz -C /mydata/
```

```shell
# 编辑环境变量配置文件
vim /etc/profile
# 在最后一行添加
export GOROOT=/mydata/go
export PATH=$PATH:$GOROOT/bin
# 刷新配置文件
source /etc/profile

go version
```

# 安装Gaea

## 安装Gaea

下载Gaea的源码，直接下载`zip`包即可，下载地址：[https://github.com/XiaoMi/Gaea](https://link.zhihu.com/?target=https%3A//github.com/XiaoMi/Gaea)

将下载好的压缩包进行解压操作，这里我们解压到`/mydata/gaea/`目录下：

```shell
unzip Gaea-master.zip
```

进入`/mydata/gaea/`目录下，使用`make`命令对源码编译：

```shell
make build
```

`注意`：由于网络问题，某些Go的依赖会下载不下来导致编译失败，多尝试几次即可成功；

编译完成后在`/mydata/gaea/bin`目录下会生成Gaea的执行文件`gaea`：

```shell
cd /mydata/gaea/bin/
```

由于我们没有搭建`etcd`配置中心，所以需要修改本地配置文件`/mydata/gaea/etc/gaea.ini`，将配置类型改为`file`：

```
; 配置类型，目前支持file/etcd两种方式，file方式不支持热加载
config_type=file
```

添加namespace配置文件，用于配置我们的主从数据库信息，配置文件地址：`/mydata/gaea/etc/file/namespace/mall_namespace_1.json`

配置文件内容如下：

```json
{
    "name": "mall_namespace_1",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "mall": true
    },
    "default_phy_dbs": {
        "mall": "mall"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "root",
            "master": "10.4.7.71:3307",
            "slaves": ["10.4.7.71:3308"],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": null,
    "users": [
        {
            "user_name": "root",
            "password": "root",
            "namespace": "mall_namespace_1",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-0",
    "global_sequences": null
}
```



## 在Docker容器中运行

由于官方只提供了Linux下直接安装运行的方式，这里我们提供另一种运行方式，在Docker容器中作为服务运行。



### 打包成Docker镜像

Docker Hub 中并没有打包好的Gaea镜像，我们需要自行构建一个，下面详细介绍下如何构建Gaea的Docker镜像。

vim Dockerfile

```dockerfile
# 该镜像需要依赖的基础镜像
FROM golang:latest
# 将当前目录下的gaea源码包复制到docker容器的/go/Gaea-master目录下，对于.tar.gz文件会自动解压
ADD Gaea-master.tar.gz /go/Gaea-master
# 将解压后的源码移动到/go/gaea目录中去
RUN bash -c 'mv /go/Gaea-master/Gaea-master /go/gaea'
# 进入/go/gaea目录
WORKDIR /go/gaea
# 将gaea源码进行打包编译
RUN go env -w GOPROXY=https://goproxy.cn
RUN bash -c 'make build'
# 声明服务运行在13306端口
EXPOSE 13306
# 指定docker容器启动时执行的命令
ENTRYPOINT ["/go/gaea/bin/gaea"]
```

在此之前我们需要把Gaea的源码压缩包转换为`.tar.gz`格式方便在Docker容器中的解压，可以使用`压缩软件`来实现：

```shell
tar zcvf Gaea-master.tar.gz gaea
```

之后使用Docker命令构建Gaea的Docker镜像：

```
docker build -t gaea:1.0.2 .
```

将本地安装的Gaea配置文件复制到`/mydata/gaea-docker/etc/`目录下：

```shell
cp -r /mydata/gaea/etc/ /mydata/gaea-docker/etc/
```

使用Docker命令启动Gaea容器：

```shell
docker run -p 13306:13306 --name gaea \
-v /mydata/gaea-docker/etc:/go/gaea/etc \
-d gaea:1.0.2
```



## 测试读写分离

> 测试思路：首先我们关闭从实例的主从复制，然后通过Gaea代理来操作数据库，插入一条数据，如果主实例中有这条数据而从实例中没有，说明写操作是走的主库。然后再通过Gaea代理查询该表数据，如果没有这条数据，表示读操作走的是从库，证明读写分离成功。

- 通过Navicat连接到Gaea代理，注意此处账号密码为Gaea的namespace中配置的内容，端口为Gaea的服务端口；

   <img src="images/image-20210912163649523.png" alt="image-20210912163649523" style="zoom:50%;" />

* 通过Navicat分别连接到主库和从库，用于查看数据，此时建立了以下三个数据库连接

    ![image-20210912214848876](images/image-20210912214848876.png)

- 通过stop slave命令关闭mysql-slave实例的主从复制功能：

  ![image-20210912214237345](images/image-20210912214237345.png)

通过Gaea代理在test表中插入一条数据：



在主库中查看test表的数据，发现已有该数据：



在从库中查看test表的数据，发现没有该数据，证明写操作走的是主库：



直接在代理中查看test表中的数据，发现没有该数据，证明读操作走的是从库。



# 配置说明

文档：https://github.com/XiaoMi/Gaea/blob/master/docs/configuration.md

namespace的配置格式为json，包含分表、非分表、实例等配置信息，都可在运行时改变。

## namespace配置说明

namespace的配置格式为json，包含分表、非分表、实例等配置信息，都可在运行时改变。namespace的配置可以直接通过web平台进行操作，使用方不需要关心json里的内容，如果有兴趣参与到gaea的开发中，可以关注下字段含义，具体解释如下,格式为字段名称、类型、内容含义。

| 字段名称             | 字段类型   | 字段含义                                                     |
| -------------------- | ---------- | ------------------------------------------------------------ |
| name                 | string     | namespace名称                                                |
| online               | bool       | 是否在线，逻辑上下线使用                                     |
| read_only            | bool       | 是否只读，namespace级别                                      |
| allowed_dbs          | map        | 数据库集合                                                   |
| default_phy_dbs      | map        | 默认数据库名, 与allowed_dbs一一对应                          |
| slow_sql_time        | string     | 慢sql时间，单位ms                                            |
| black_sql            | string数组 | 黑名单sql                                                    |
| allowed_ip           | string数组 | 白名单IP                                                     |
| slices               | map数组    | 一主多从的物理实例，slice里map的具体字段可参照slice配置      |
| shard_rules          | map数组    | 分库、分表、特殊表的配置内容，具体字段可参照shard配置        |
| users                | map数组    | 应用端连接gaea所需要的用户配置，具体字段可参照users配置      |
| global_sequences     | map        | 生成全局唯一序列号的配置, 具体字段可参考全局序列号配置       |
| default_slice        | string     | show语句默认的执行分片                                       |
| open_general_log     | bool       | 是否开启审计日志, [如何开启](https://github.com/XiaoMi/Gaea/issues/109) |
| max_sql_execute_time | int        | 应用端查询最大执行时间, 超时后会被自动kill, 为0默认不开启此功能 |
| max_sql_result_size  | int        | gaea从后端mysql接收结果集的最大值, 限制单分片查询行数, 默认值10000, -1表示不开启 |

### slice配置

| 字段名称         | 字段类型   | 字段含义                                       |
| ---------------- | ---------- | ---------------------------------------------- |
| name             | string     | 分片名称，自动、有序生成                       |
| user_name        | string     | 连接后端mysql所需要的用户名称                  |
| password         | string     | 连接后端mysql所需要的用户密码                  |
| master           | string     | 主实例地址                                     |
| slaves           | string数组 | 从实例地址列表                                 |
| statistic_slaves | string数组 | 统计型从实例地址列表                           |
| capacity         | int        | gaea_proxy与每个实例的连接池大小               |
| max_capacity     | int        | gaea_proxy与每个实例的连接池最大大小           |
| idle_timeout     | int        | gaea_proxy与后端mysql空闲连接存活时间，单位:秒 |

### shard配置

这里列出了一些基本配置参数, 详细配置请参考[分片表配置](https://github.com/XiaoMi/Gaea/blob/master/docs/shard.md)

如需要了解每种规则详细库表对照示例，可以查看[分片规则示例说明](https://github.com/XiaoMi/Gaea/blob/master/docs/shard-example.md)

| 字段名称  | 字段类型 | 字段含义                  |
| --------- | -------- | ------------------------- |
| db        | string   | 分片表所在DB              |
| table     | string   | 分片表名                  |
| type      | string   | 分片类型                  |
| key       | string   | 分片列名                  |
| locations | list     | 每个slice上分布的分片个数 |
| slices    | list     | slice列表                 |
| databases | list     | mycat分片规则后端实际DB名 |

### users配置

| 字段名称       | 字段类型 | 字段含义                                             |
| -------------- | -------- | ---------------------------------------------------- |
| user_name      | string   | 用户名                                               |
| password       | string   | 用户密码                                             |
| namespace      | string   | 对应的命名空间                                       |
| rw_flag        | int      | 读写标识, 只读=1, 读写=2                             |
| rw_split       | int      | 是否读写分离, 非读写分离=0, 读写分离=1               |
| other_property | int      | 目前用来标识是否走统计从实例, 普通用户=0, 统计用户=1 |

### 全局序列号配置

| 字段名称   | 字段类型 | 字段含义                                             |
| ---------- | -------- | ---------------------------------------------------- |
| db         | string   | 使用全局序列号的表所在的db的逻辑db名                 |
| table      | string   | 使用全局序列号的表的逻辑表名                         |
| type       | string   | 目前只支持mycat方式                                  |
| pk_name    | string   | 使用全局序列号的列名，单表只允许一个列使用全局序列号 |
| slice_name | string   | mycat_sequence表所在分片                             |

## 配置示例

```
{
    "name": "gaea_namespace_1",
    "online": true,
    "read_only": true,
    "allowed_dbs": {
        "db_ks": true,
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_ks": "db_ks",
        "db_mycat": "db_mycat_0"
    }, 
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null, 
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:3306",
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:3307",
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        }
    ],
    "shard_rules": [
        {
            "db": "db_ks",
            "table": "tbl_ks",
            "type": "hash",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_child",
            "type": "linked",
            "key": "id",
            "parent_table": "tbl_ks"
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_global",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_range",
            "type": "range",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "table_row_limit": 100
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_year",
            "type": "date_year",
            "key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "date_range": [
                "2014-2017",
                "2018-2019"
            ]
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_month",
            "type": "date_month",
            "key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "date_range": [
                "201405-201406",
                "201408-201409"
            ]
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_day",
            "type": "date_day",
            "key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "date_range": [
                "20140901-20140905",
                "20140907-20140908"
            ]
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "mycat_mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_child",
            "type": "linked",
            "parent_table": "tbl_mycat",
            "key": "id"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_murmur",
            "type": "mycat_murmur",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_0",
                "db_mycat_1",
                "db_mycat_2",
                "db_mycat_3"
            ],
            "seed": "0",
            "virtual_bucket_times": "160"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_long",
            "type": "mycat_long",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
            "partition_count": "4",
            "partition_length": "256"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_global",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_string",
            "type": "mycat_string",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
            "partition_count": "4",
            "partition_length": "256",
            "hash_slice": "20"
        }
    ],
    "global_sequences": [
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "test",
            "pk_name": "id"
        },
        {
            "db": "db_ks",
            "table": "tbl_ks",
            "type": "test",
            "pk_name": "user_id"
        }
    ],
    "users": [
        {
            "user_name": "test_shard",
            "password": "test_shard",
            "namespace": "gaea_namespace_1",
            "rw_flag": 2,
            "rw_split": 1
        }
    ], 
    "default_slice": "slice-0",
    "open_general_log": false,
    "max_sql_execute_time": 5000,
    "max_sql_result_size": 10
}
```

本配置截取自proxy/plan/plan_test.go, 如果对Gaea分表有困惑, 也可以参考这个包下的测试用例. 下面将结合该配置示例介绍Gaea的namespace配置细节.

namespace名称为`gaea_namespace_1`. 在该namespace的`users`字段中添加一个gaea用户`test_shard`. 特别注意Gaea中的`用户名+密码`是全局唯一的 (映射到唯一的namespace). 该用户是读写用户, 且使用读写分离.

在namespace中通过`allowed_dbs`字段配置了两个可用的数据库, 另一个相关的字段为`default_phy_dbs`, 该字段仅用于mycat分库路由的场景, 用于标记后端实际库名. 如果没有使用mycat路由, 则可以只配置`allowed_dbs`字段, 不配置`default_phy_dbs`字段.

通过`slices`字段配置后端的slice. 一个slice实际上对应着一组MySQL实例, 可以包含一主多从. slice的名称目前必须使用`slice-0`, `slice-1`这样的格式, 如果自定义slice名称会出现找不到默认slice的问题.

在`shard_rules`字段中配置分片表信息. 按照Gaea处理方式, 可以将分片表分为3类: kingshard路由模式的分片表, mycat路由模式的分片表, 全局表.



### kingshard路由

kingshard路由模式下, 分片表要求后端数据库的库名相同, 子表的表名为`table_后缀`的模式.

```
{
    "db": "db_ks",
    "table": "tbl_ks",
    "type": "hash",
    "key": "id",
    "locations": [
        2,
        2
    ],
    "slices": [
        "slice-0",
        "slice-1"
    ]
}
```

以这个kingshard hash分片表配置为例, 路由规则为hash, 逻辑表名为tbl_ks, locations 2,2表示有两个slice, 每个slice上面分配两张子表, `slices`配置了两个slice的名称. 那么后端数据库的子表需要按照以下规则创建:

| slice   | db    | table       |
| ------- | ----- | ----------- |
| slice-0 | db_ks | tbl_ks_0000 |
| slice-0 | db_ks | tbl_ks_0001 |
| slice-1 | db_ks | tbl_ks_0002 |
| slice-1 | db_ks | tbl_ks_0003 |

其他kingshard路由的表名映射关系均类似, 再以range路由举例:

```
{
    "db": "db_ks",
    "table": "tbl_ks_month",
    "type": "date_month",
    "key": "create_time",
    "slices": [
        "slice-0",
        "slice-1"
    ],
    "date_range": [
        "201405-201406",
        "201408-201409"
    ]
}
```

| slice   | db    | table         |
| ------- | ----- | ------------- |
| slice-0 | db_ks | tbl_ks_201405 |
| slice-0 | db_ks | tbl_ks_201406 |
| slice-1 | db_ks | tbl_ks_201408 |
| slice-1 | db_ks | tbl_ks_201409 |

kingshard路由不需要配置`databases`字段, 因为后端数据库名与逻辑库名相同.



### mycat路由

mycat路由与kingshard不完全相同, Gaea主要兼容了mycat的分库路由模式.

```
{
    "db": "db_mycat",
    "table": "tbl_mycat_murmur",
    "type": "mycat_murmur",
    "key": "id",
    "locations": [
        2,
        2
    ],
    "slices": [
        "slice-0",
        "slice-1"
    ],
    "databases": [
        "db_mycat_0",
        "db_mycat_1",
        "db_mycat_2",
        "db_mycat_3"
    ],
    "seed": "0",
    "virtual_bucket_times": "160"
}
```

| slice   | db         | table            |
| ------- | ---------- | ---------------- |
| slice-0 | db_mycat_0 | tbl_mycat_murmur |
| slice-0 | db_mycat_1 | tbl_mycat_murmur |
| slice-1 | db_mycat_2 | tbl_mycat_murmur |
| slice-1 | db_mycat_3 | tbl_mycat_murmur |

其中`databases`字段需要按路由顺序指定后端数据库的实际库名, 且数量需要与`locations`的总和相等.

### 全局表路由

全局表路由与mycat路由配置类似, 但是可以不指定`databases`. 如果不指定, 则全局表在各个后端的数据库名和表名均相同.