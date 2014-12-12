package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var amqpUri string
var connection *amqp.Connection
var channel *amqp.Channel

func init_amqp_publisher() {
	amqp_config := get_amqp_config()
	amqpUri = fmt.Sprintf("amqp://%s:%s@%s:%d/", amqp_config.UserName, amqp_config.Password, amqp_config.Host, amqp_config.Port)
	var err error
	connection, err = amqp.Dial(amqpUri)
	if err == nil {
		channel, err = connection.Channel()
	}

	go func() {
		var closedConnChannel = connection.NotifyClose(make(chan *amqp.Error))
		closeErr := <-closedConnChannel
		log.Printf("connection error %s,  attempting reconnect after 10 sec", closeErr)
		time.Sleep(5 * time.Second)
		init_amqp_publisher()
	}()
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
