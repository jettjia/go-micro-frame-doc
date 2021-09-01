package consumer

import (
	"github.com/houseofcat/turbocookedrabbit/v2/pkg/tcr"
	"log"
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

	rabbitService.GetConsumerConfig("turboQueueName")
}
