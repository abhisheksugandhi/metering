// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
//
package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"util"
)

var analytics_client util.AnalyticsClient = util.AnalyticsClient{}
var dm_client util.GDMClient = util.GDMClient{}

func init_amqp_subscriber() {
	analytics_client.Init("http://localhost:8085/printEvents")
	amqp_config := get_amqp_config()
	amqp_uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", amqp_config.UserName, amqp_config.Password, amqp_config.Host, amqp_config.Port)
	dm_client.Init("http://localhost:8085/test_distance", amqp_uri)

	c, err := subscribe_to_queues(amqp_uri)
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("running forever")
	select {}

	log.Printf("shutting down")
	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func subscribe_to_queues(amqpURI string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     "",
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	queueName := "elroy.analytics.v0"
	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, queue.Name)

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	analytic_deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle_analytics(analytic_deliveries, c.done)

	queueName = "elroy.eta.v0"
	queue, err = c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, queue.Name)

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	eta_deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle_eta(eta_deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle_analytics(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf("got event of length %d: %q", len(d.Body), d.Body)
		analytics_client.Call(d.Body)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func handle_eta(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf("got event of length %d: %q", len(d.Body), d.Body)
		dm_client.GetDM(d.Body, d.ReplyTo, d.MessageId, d.CorrelationId)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
