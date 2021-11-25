konga后台使用介绍

Kong 官网：https://konghq.com/
Kong GitHub地址：https://github.com/kong/kong
konga GitHub地址：https://github.com/pantsel/konga

# 安装

## 安装kong

参考 readme 中 安装kong



## 安装konga

方式1，开发环境

```shell
docker run -d -p 1337:1337 --name konga pantsel/konga
```

方式2，生产环境

安装数据库

```shell
docker run --network=kong-net --rm pantsel/konga -c prepare -a postgres -u postgresql://konga:konga@kong-database/konga

```

```shell
docker run -d -p 1337:1337 --name kong-konga \ 
--network=kong-net \ 
-e "TOKEN_SECRET=${自己随机生成一个}" \ 
-e "DB_ADAPTER=postgres" \ 
-e "DB_HOST=kong-database" \ 
-e "DB_PORT=5432" \ 
-e "DB_USER=konga" \ 
-e "DB_PASSWORD=konga" \ 
-e "DB_DATABASE=konga" \ 
-e "NODE_ENV=development" \ 
pantsel/konga

```

```awk
http://<your-servers-public-ip-or-host>:1337
```



# Konga 连接 Kong

 ![image-20211125170417188](assets/image-20211125170417188.png)

# Dashboard

 ![image-20211125170504830](assets/image-20211125170504830.png)

仪表板显示有关你当前连接的 Kong 实例，基础数据库和可用插件的基本信息。更多详细信息可在 INFO 页面中找到。

# Snapshots

快照功能允许你轻松地跨节点备份，恢复和移动 Kong 配置。你还可以安排 Kong 实例的自动快照。

 ![image-20211125170628095](assets/image-20211125170628095.png)

# Settings

设置页面提供了一种配置Konga并为用户帐户设置基本ACL的简单方法。请记住，用户权限是全局设置的，并将用户帐户称为实体。尚不支持单个用户ACL。

 ![image-20211125170733832](assets/image-20211125170733832.png)

# 配置路由转发- todo

​		service服务是一个抽象服务层，可以用于指向具体物理服务（target），也可以指向upstream用于实现物理服务的负载效果。一个service对于upstream、target都是一对一的关系。

​		routes路由匹配客户端的请求规则。匹配成功后分配到service层。一个路由指向一个service，一个service可以被多个不通规则的路由routes指向。相当于nginx配置中的 location

​		upstream主要用于实现kong的负载功能，一个service匹配到一个upstream后，upstream可以指向多个target服务以此来实现负载的效果。

​		target服务端指的就是最后的物理服务，当然也可以是其他的虚拟服务。



## service

| 属性                    | 描述                                                         |
| ----------------------- | ------------------------------------------------------------ |
| name（必填）            | 服务名称.                                                    |
| tags（可选）            | 可选地向服务添加标记                                         |
| url（可选）             | 将协议、主机、端口和路径立即设置成简短的属性。这个属性是只写的（管理API从来不“返回”url） |
| protocol（必填）        | 该协议用于与upstream通信。它可以是http（默认）或https。      |
| host（必填）            | upstream服务器的主机。                                       |
| port（必填）            | upstream服务器端口。默认为80                                 |
| path（可选）            | 在向upstream服务器请求中使用的路径。默认为空。               |
| retries（可选）         | 在代理失败的情况下执行的重试次数。默认值是5。                |
| connect_timeout（可选） | 建立到upstream服务器的连接的超时时间。默认为60000。          |
| write_timeout（可选）   | 将请求发送到upstream服务器的两个连续写操作之间的超时时间。默认为60000。 |
| read_timeout（可选）    | 将请求发送到upstream服务器的两个连续读取操作之间的超时时间。默认为60000。 |



## router

| name(可选)                        | 定义名称                                                     |
| --------------------------------- | ------------------------------------------------------------ |
| tags(可选)                        | 向路由添加标记                                               |
| protocols(可选)                   | 这条路线应该允许的协议列表。默认情况下，它是“http”、“https”，这意味着路由接受这两种方式。当设置为“HTTPS”时，HTTP请求会被请求升级到HTTPS。通过表单编码，符号是协议=http&协议=https。使用JSON，使用数组。 |
| methods(半可选)                   | 与此路由相匹配的HTTP方法列表。例如: ["GET", "POST"].至少有一个主机、路径或方法必须被设置。用表单编码参数是methods[]=GET&methods[]=OPTIONS。使用JSON，使用数组。 |
| hosts(半可选)                     | 与此路径匹配的域名列表。例如:example.com. 至少有一个主机、路径或方法必须被设置。用表单编码参数是 hosts[]=foo.com&hosts[]=bar.com。使用JSON，使用数组。 |
| paths(半可选)                     | 与此路径相匹配的路径列表。例如: /my-path.至少有一个主机、路径或方法必须被设置。用表单编码参数是 paths[]=/foo&paths[]=/bar. 使用JSON，使用数组。 |
| Https redirect status code (可选) | 当路由的所有属性都匹配（协议除外）时，即如果请求的协议是HTTP而不是HTTPS，则状态代码Kong将响应。如果字段设置为301、302、307或308，则位置头由Kong注入。默认为426。 |
| Regex priority(可选)              | 当多个路由同时使用正则表达式匹配某个给定请求时，用来选择哪个路由解析该请求的数字。当两个路由匹配路径并且具有相同的regex_优先级时，将使用较旧的路由（最低创建位置）。注意，非regex路由的优先级不同（较长的非regex路由在较短的路由之前匹配）。默认为0。 |
| strip_path(可选)                  | 当通过一条路径匹配一条路径时，从upstream请求URL中剥离匹配的前缀。默认值为true。 |
| preserve_host(可选)               | 当通过一个主机域名匹配一条路由时，在upstream请求头中使用请求主机头。默认设置为false，upstream主机头将是服务主机的主机头。 |



# 集成consul实现服务发现和负载均衡-todo

 查看info中，是否有成功配置好 dns_server信息；

<img src="assets/image-20210825114630590.png" alt="image-20210825114630590" style="zoom:50%;" />

10.4.7.71:8000/g/v1/goods

# 配置jwt实现登录校验- todo

1) 创建 consumers

2) 为这个 consumer 配置 jwt

3) 记住内容中的 key

4) 配置全局的 plugin

4.1 ) 设置全局的token plugin -> add plugin -> jwt -> header names : x-token

注意：这里的 key, secret 要和代码中已经存在的 Jwt 配置保持一致。否则kong的jwt通过认证了，代码中的jwt还没有通过；

实际开发中，可以只用 kong中的jwt；然后用（4.2）中的方式，生成token；也可以写一个接口给到前端调取生成token

 <img src="assets/image-20210825115737779.png" alt="image-20210825115737779" style="zoom:33%;" />

4.2) jwt.io生成 token

 <img src="assets/image-20210825115310510.png" alt="image-20210825115310510" style="zoom:33%;" />

 <img src="assets/image-20210825115348955.png" alt="image-20210825115348955" style="zoom:33%;" />

上面这个 secret 从 konga的管理后台里选择；

# 配置反爬和ip黑名单

plugins -> security

ip 限制

plugins -> security->ip restriction

 <img src="assets/image-20210825124913264.png" alt="image-20210825124913264" style="zoom:33%;" />



ip反爬

plugins -> security->bot detection

 <img src="assets/image-20210825124805220.png" alt="image-20210825124805220" style="zoom:33%;" />



# 配置熔断和限流

## 请求大小限制

### 启用大小限制

1、在服务上启用插件

```kotlin
$ curl -X POST http://kong:8001/services/{service}/plugins \
    --data "name=request-size-limiting"  \
    --data "config.allowed_payload_size=128"
复制代码
```

2、在路由上启用插件

```kotlin
$ curl -X POST http://kong:8001/routes/{route_id}/plugins \
    --data "name=request-size-limiting"  \
    --data "config.allowed_payload_size=128"
复制代码
```

3、在消费者中启用插件

```kotlin
$ curl -X POST http://kong:8001/plugins \
    --data "name=request-size-limiting" \
    --data "consumer_id={consumer_id}"  \
    --data "config.allowed_payload_size=128"
复制代码
curl -X POST http://172.16.60.8:8001/routes/b0db420a-d3c3-45ee-8b25-11f3fd8ca283/plugins \
    --data "name=request-size-limiting"  \
    --data "config.allowed_payload_size=-20"
    
{
	"created_at": 1609212809,
	"id": "08112606-e023-4d1b-ba76-853832585feb",
	"tags": null,
	"enabled": true,
	"protocols": ["grpc", "grpcs", "http", "https"],
	"name": "request-size-limiting",
	"consumer": null,
	"service": null,
	"route": {
		"id": "b0db420a-d3c3-45ee-8b25-11f3fd8ca283"
	},
	"config": {
		"size_unit": "megabytes",
		"allowed_payload_size": -20
	}
}
```

### 访问测试

 ![image-20211125173429561](assets/image-20211125173429561.png)





## 请求速率限制

### 启用速率限制

速率限制开发人员在给定的时间、分钟、小时、日、月或年的时间内可以发出多少HTTP请求

1、在服务上启用插件

```kotlin
$ curl -X POST http://kong:8001/services/{service}/plugins \
    --data "name=rate-limiting"  \
    --data "config.second=5" \
    --data "config.hour=10000"
复制代码
```

2、在路由上启用插件

```kotlin
$ curl -X POST http://kong:8001/routes/{route_id}/plugins \
    --data "name=rate-limiting"  \
    --data "config.second=5" \
    --data "config.hour=10000"
复制代码
```

3、在消费者中启用插件

```kotlin
$ curl -X POST http://kong:8001/plugins \
    --data "name=rate-limiting" \
    --data "consumer_id={consumer_id}"  \
    --data "config.second=5" \
    --data "config.hour=10000"
复制代码
curl -X POST http://172.16.60.8:8001/routes/b0db420a-d3c3-45ee-8b25-11f3fd8ca283/plugins \
    --data "name=rate-limiting"  \
    --data "config.second=2" \
    --data "config.hour=10000"

{
	"created_at": 1609213117,
	"id": "0329fd21-3844-4a84-b325-68282577c78a",
	"tags": null,
	"enabled": true,
	"protocols": ["grpc", "grpcs", "http", "https"],
	"name": "rate-limiting",
	"consumer": null,
	"service": null,
	"route": {
		"id": "b0db420a-d3c3-45ee-8b25-11f3fd8ca283"
	},
	"config": {
		"minute": null,
		"redis_host": null,
		"redis_timeout": 2000,
		"limit_by": "consumer",
		"hour": 10000,
		"policy": "cluster",
		"month": null,
		"redis_password": null,
		"second": 2,
		"day": null,
		"hide_client_headers": false,
		"path": null,
		"redis_database": 0,
		"year": null,
		"redis_port": 6379,
		"header_name": null,
		"fault_tolerant": true
	}
}
```


### 测试

 ![image-20211125173545345](assets/image-20211125173545345.png)





## 请求终止

这个插件使用指定的状态码和消息终止传入的请求。这允许(暂时)停止服务或路由(或废弃的API实体)上的通信，甚至阻塞用户

### 启用请求终止插件

1、在服务上启用插件

```kotlin
$ curl -X POST http://kong:8001/services/{service}/plugins \
    --data "name=request-termination"  \
    --data "config.status_code=403" \
    --data "config.message=So long and thanks for all the fish!"
复制代码
```

2、在路由上启用插件

```kotlin
$ curl -X POST http://kong:8001/routes/{route_id}/plugins \
    --data "name=request-termination"  \
    --data "config.status_code=403" \
    --data "config.message=So long and thanks for all the fish!"
复制代码
curl -X POST http://172.16.60.8:8001/routes/b0db420a-d3c3-45ee-8b25-11f3fd8ca283/plugins \
    --data "name=request-termination"  \
    --data "config.status_code=403" \
    --data "config.message=The service is under maintenance"

{
	"created_at": 1609213373,
	"id": "77c99041-0a99-4388-b4fa-9378791c3839",
	"tags": null,
	"enabled": true,
	"protocols": ["grpc", "grpcs", "http", "https"],
	"name": "request-termination",
	"consumer": null,
	"service": null,
	"route": {
		"id": "b0db420a-d3c3-45ee-8b25-11f3fd8ca283"
	},
	"config": {
		"status_code": 403,
		"content_type": null,
		"body": null,
		"message": "The service is under maintenance"
	}
}
```


### 测试

 ![image-20211125173645669](assets/image-20211125173645669.png)

### 备注

config.status_code ：要发送的响应代码。

config.message ：如果使用默认响应生成器，则发送的消息。

config.body ： 要发送的raw响应体，跟`config.message`相互排斥

config.content_type : 为`config.body`配置raw响应的内容类型，默认：`application/json; charset=utf-8`



# Nginx配置 VS Kong配置

我们来看一个典型的Nginx的配置对应在Kong上是怎么样的，下面是一个典型的Nginx配置

```nginx
upstream passportUpstream {
        server localhost:8080 weight=100;
}
server {
        listen 80;
        location /hello {
        	proxy_pass http://passportUpstream;
        }
}
```

下面我们来看看其对应Kong中的配置

```shell
# 配置 upstream 
curl -X POST http://localhost:8001/upstreams 
    --data "name=passportUpstream" 

# 配置 target 
curl -X POST http://localhost:8001/upstreams/passport/targets 
    --data "target=localhost:8080" --data "weight=100" 

# 配置 service 
curl -X POST http://localhost:8001/services 
    --data "name=getUserInfo" --data "host=passportUpstream" 

# 配置 route 
curl -X POST http://localhost:8001/routes 
    --data "paths[]=/user"
    --data "service.id=8695cc65-16c1-43b1-95a1-5d30d0a50409" 

curl -X POST http://localhost:8001/routes 
    --data "hosts[]=*.example.com,test.com,*.abc.com" 
    --data "service.id=8695cc65-16c1-43b1-95a1-5d30d0a50409" 
```

这一切配置都是通过其Http Restful API 来动态实现的，无需我们在手动的 reload Nginx.conf 。开发的同学看到这是不是感觉到很幸福了。

在上述的配置中涉及到了几个概念：upstrean、target、service、route等概念，它们是Kong的几个核心概念，也是我们在使用Kong Api 时经常打交道的，下面我们就其几个核心概念做一下简单的说明。



# kong关键词和名词解析

**Upstream**
Upstream 对象表示虚拟主机名，可用于通过多个服务（目标）对传入请求进行负载均衡。例如：service.v1.xyz 为Service对象命名的上游host是service.v1.xyz对此服务的请求将代理到上游定义的目标。

**Target**
目标IP地址/主机名，其端口表示后端服务的实例。每个上游都可以有多个target,并且可以动态添加Target。

由于上游维护Target的更改历史记录，因此无法删除或者修改Target。要禁用目标，请发布一个新的Targer weight=0,或者使用DELETE来完成相同的操作。

**Service**
顾名思义，服务实体是每个上游服务的抽象。服务的示例是数据转换微服务，计费API等。

服务的主要属性是它的URL（其中，Kong应该代理流量），其可以被设置为单个串或通过指定其protocol， host，port和path。

服务与路由相关联（服务可以有许多与之关联的路由）。路由是Kong的入口点，并定义匹配客户端请求的规则。一旦匹配路由，Kong就会将请求代理到其关联的服务。

**Route**
路由实体定义规则以匹配客户端的请求。每个Route与一个Service相关联，一个服务可能有多个与之关联的路由。与给定路由匹配的每个请求都将代理到其关联的Service上。可以配置的字段有

- hosts
- paths
- methods

Service 和 Route 的组合（以及它们之间的关注点分离）提供了一种强大的路由机制，通过它可以在Kong中定义细粒度的入口点，从而使基础架构路由到不同上游服务。

**Consumer**
Consumer 对象表示服务的使用者或者用户。您可以依靠Kong作为主数据库存储，也可以将使用者列表与数据库映射，以保持Kong与现有的主数据存储之间的一致性。

**Plugin**
插件实体表示将在HTTP请求/响应生命周期期间执行的插件配置。它是如何为在Kong后面运行的服务添加功能的，例如身份验证或速率限制。

将插件配置添加到服务时，客户端向该服务发出的每个请求都将运行所述插件。如果某个特定消费者需要将插件调整为不同的值，您可以通过创建一个单独的插件实例，通过service和consumer字段指定服务和消费者 。

**对应关系**
Upstream : target -> 1:n
Service:Upstream -> 1:1 or 1:0 (Service 可以直接指向具体的Target，相当于不做负载均衡)
Service : Route -> 1:n

**Note:** Client请求的流量通过Route指向与之相关的Service，如果配置插件的话就会作用插件，Service接到流量后给相应的Upstream的服务上面。