package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/houseofcat/turbocookedrabbit/v2/pkg/tcr"
	"log"
	"time"
)

func main() {

	config, err := tcr.ConvertJSONFileToConfig("16-amqp/rabbitmq/00-rabbitmq/config.json") // Load Configuration On Startup
	if err != nil {
		log.Fatal(err)
	}

	rabbitService, err := tcr.NewRabbitService(config, "", "", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	// Then publish (this time with a confirmation/context)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rabbitService.Publisher.PublishWithConfirmationContext(
		ctx,
		&tcr.Letter{
			LetterID: id,
			Body:     []byte("Hello World"),
			Envelope: &tcr.Envelope{
				Exchange:     "MyDeclaredExchangeName",
				RoutingKey:   "MyDeclaredQueueName",
				ContentType:  "text/plain",
				Mandatory:    false,
				Immediate:    false,
				DeliveryMode: 2,
			},
		})

	fmt.Println("create success")

}
