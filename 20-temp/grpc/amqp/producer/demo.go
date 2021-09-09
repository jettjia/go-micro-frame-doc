package producer

import (
	"fmt"
	"log"

	"go-micro-module/16-amqp/rabbitmq/01-gin-web/utils/amqpRabbit"
	"go-micro-module/20-temp/grpc/global"
)

// 发送消息到 mq
func TestGoods(body string) {
	exchangeName := "test-exchange"
	exchangeType := "direct"
	routingKey := "test-key"
	reliable := true

	c := global.ServerConfig.MqInfo
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port)

	if err := amqpRabbit.Publish(uri, exchangeName, exchangeType, routingKey, body, reliable); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("发送的数据内容是 %dB OK", len(body))
}
