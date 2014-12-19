// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
//
package util

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

//var analytics_client util.AnalyticsClient = util.AnalyticsClient{}
//var dm_client util.GDMClient = util.GDMClient{}

type AmqpSubscriber struct {
	amqp_uri string
	queueName string
	handlerFunction func(<-chan amqp.Delivery, chan error)
	consumer *Consumer
}

func (subscriber *AmqpSubscriber) Init(amqpUri string, queueName string, handlerFunction func(<-chan amqp.Delivery, chan error)) {
	subscriber.queueName = queueName
	subscriber.handlerFunction = handlerFunction
	subscriber.amqp_uri = amqpUri
	subscriber.init_internal()
} 

func (subscriber *AmqpSubscriber) init_internal() {
	//dm_client.Init("http://localhost:8085/test_distance", amqp_uri)

	c, err := subscribe_to_queues(subscriber.amqp_uri, subscriber.queueName, subscriber.handlerFunction)
	if err != nil {
		log.Fatalf("%s", err)
	}
	subscriber.consumer = c

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

func subscribe_to_queues(amqpURI string, queueName string, handlerFunction func(<-chan amqp.Delivery, chan error)) (*Consumer, error) {
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
		log.Printf("got error while dialing %s", err)
	}

	go func() {
		var closedConnChannel = c.conn.NotifyClose(make(chan *amqp.Error))
		closeErr := <-closedConnChannel
		log.Printf("connection error %s,  attempting reconnect after 5 sec", closeErr)
		time.Sleep(5 * time.Second)
		subscribe_to_queues(amqpURI, queueName, handlerFunction)
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

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

	go handlerFunction(analytic_deliveries, c.done)

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
		d.Ack(false)
		log.Printf("got event of length %d: %q", len(d.Body), d.Body)
		//analytics_client.Call(d.Body)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func handle_eta(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		d.Ack(false)
		log.Printf("got event of length %d: %q", len(d.Body), d.Body)
		//dm_client.GetDM(d.Body, d.ReplyTo, d.MessageId, d.CorrelationId)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
