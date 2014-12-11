package util

import (
	"fmt"
	"github.com/streadway/amqp"
)

type AmqpPublisher struct {
	AmqpUri    string
	Connection *amqp.Connection
}

func (publisher *AmqpPublisher) Init(amqpUri string) {
	connection, err := amqp.Dial(amqpUri)
	if err != nil {
		panic(err)
	}
	publisher.Connection = connection
}

func (publisher *AmqpPublisher) Publish_Amqp(message []byte, routingKey string, messageId string, correlationId string) error {
	channel, err1 := publisher.Connection.Channel()
	if err1 != nil {
		panic(err1)
	}
	exchange := ""
	fmt.Printf("publishing to queue %s message of length %d\n", routingKey, len(message))
	err := channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		true,  // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          message,
			MessageId:     messageId,
			CorrelationId: correlationId,
		})
	if err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}
	fmt.Println("done publishing")
	return err
}
