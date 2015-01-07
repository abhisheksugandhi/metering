package main

import (
	"flag"
	"flywheel/util"
	"fmt"
	"log"
	"os"
)

func run_distance_matrix_client_service() {
	config := util.Config{}
	config.Init()
	amqp_config := config.AmqpConfig
	amqp_uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", amqp_config.UserName, amqp_config.Password, amqp_config.Host, amqp_config.Port)

	dm_client := GDMClient{}
	dm_client.Init("http://localhost:8085/test_distance", amqp_uri)
	dm_client.subscribe_to_queue(amqp_uri, "elroy.eta.v0")
}

func run_analytics_client_service() {
	config := util.Config{}
	config.Init()
	amqp_config := config.AmqpConfig
	amqp_uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", amqp_config.UserName, amqp_config.Password, amqp_config.Host, amqp_config.Port)
	analytics_client := AnalyticsClient{}
	analytics_client.Init()
	amqp_subscriber := util.AmqpSubscriber{}
	amqp_subscriber.Init(amqp_uri, "elroy.analytics.v0", analytics_client.handle_analytics)
}

func main() {
	file, err := os.OpenFile("logs/metering.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		panic(err)
	}

	methodPtr := flag.String("method", "analytics_client", "method name")
	flag.Parse()
	log.Println("method:", *methodPtr)
	var method = *methodPtr
	if method == "distance_matrix_client" {
		run_distance_matrix_client_service()
	} else if method == "analytics_client" {
		run_analytics_client_service()
	} else {
		log.Printf("please run using valid -method flag")
	}
}
