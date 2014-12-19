package util

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type AmqpPublisher struct {
	AmqpUri    string
	Connection *amqp.Connection
}

func (publisher *AmqpPublisher) Init(amqpUri string) {
	publisher.AmqpUri = amqpUri
	connection, err := amqp.Dial(amqpUri)
	if err == nil {
		publisher.Connection = connection
	}

	go func() {
		var closedConnChannel = connection.NotifyClose(make(chan *amqp.Error))
		closeErr := <-closedConnChannel
		fmt.Printf("connection error %s,  attempting reconnect after 1 sec", closeErr)
		time.Sleep(5 * time.Second)
		publisher.Init(amqpUri)
	}()

}

func (publisher *AmqpPublisher) Publish_Amqp(message []byte, routingKey string, messageId string, correlationId string) error {
	channel, err1 := publisher.Connection.Channel()
	if err1 != nil {
		fmt.Printf("got error while getting channel %s", err1)
		return err1;
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
