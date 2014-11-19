package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

var amqpUri string
var connection *amqp.Connection
var channel *amqp.Channel

func init_amqp_publisher() {
	amqp_config := get_amqp_config()
	amqpUri = fmt.Sprintf("amqp://%s:%s@%s:%d/", amqp_config.UserName, amqp_config.Password, amqp_config.Host, amqp_config.Port)
	var err error
	connection, err = amqp.Dial(amqpUri)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err = connection.Channel()
	if err != nil {
		panic(err)
	}
}

func publish_amqp(routingKey string, message []byte) error {
	exchange := ""
	var err error
	if err = channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            message,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}
