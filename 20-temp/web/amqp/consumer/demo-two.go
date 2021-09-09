package consumer
//
//import (
//	"log"
//	"time"
//
//	"go-micro-module/20-temp/web/utils/amqpRabbit"
//)
//
//func init()  {
//	TestTwo()
//}
//
//func TestTwo() {
//	exchangeName := "test-exchange"
//	exchangeType := "direct"
//	bindingKey := "test-key"
//	consumerTag := "simple-consumer"
//	uri := "amqp://admin:123456@10.4.7.71:5672/"
//	queue := "test-queue"
//	lifetime := time.Duration(0) * time.Second
//
//	go func() {
//		c, err := amqpRabbit.NewConsumer(uri, exchangeName, exchangeType, queue, bindingKey, consumerTag)
//		if err != nil {
//			log.Fatalf("%s", err)
//		}
//
//		if lifetime > 0 {
//			log.Printf("running for %s", lifetime)
//			time.Sleep(lifetime)
//		} else {
//			log.Printf("running forever")
//			select {}
//		}
//
//		log.Printf("shutting down")
//
//		if err := c.Shutdown(); err != nil {
//			log.Fatalf("error during shutdown: %s", err)
//		}
//	}()
//}
