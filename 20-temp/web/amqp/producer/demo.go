package producer

import (
	"log"

	"go-micro-module/20-temp/web/utils/amqpRabbit"
)

// 发送消息到 mq
func TestGoods(body string) {
	exchangeName := "test-exchange"
	exchangeType := "direct"
	routingKey := "test-key"
	reliable := true
	uri := "amqp://admin:123456@10.4.7.71:5672/"

	if err := amqpRabbit.Publish(uri, exchangeName, exchangeType, routingKey, body, reliable); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("发送的数据内容是 %dB OK", len(body))
}
