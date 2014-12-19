package main

import (
	"bytes"
	"log"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/streadway/amqp"
	"fmt"
	"flywheel/util"
)

type AnalyticsClient struct {
	url string
	client http.Client
}

func (a *AnalyticsClient) Init(serverUrl string) {
	a.url = serverUrl
	a.client = http.Client{}
}

func (a *AnalyticsClient) Call(data []byte) {
	postData := url.Values{}
	postData.Set("content", bytes.NewBuffer(data).String())
	req, err := http.NewRequest("POST", a.url, bytes.NewBufferString(postData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := a.client.Do(req)
	if err != nil {
		fmt.Errorf("Server call error %s", err)
	}
	defer resp.Body.Close()
	log.Printf("got response status %s on analytics call\n", resp.Status)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	log.Printf("got response from GDM server as %s \n", bytes.NewBuffer(responseBody).String())
}

func (a *AnalyticsClient) handle_analytics(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		d.Ack(false)
		log.Printf("got event of length %d: %q", len(d.Body), d.Body)
		a.Call(d.Body)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func main_1() {
	config := util.Config{}
	config.Init()
	amqp_config := config.AmqpConfig
	amqp_uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", amqp_config.UserName, amqp_config.Password, amqp_config.Host, amqp_config.Port)
	analytics_client := AnalyticsClient{}
	analytics_client.Init("http://localhost:8085/printEvents")
	amqp_subscriber := util.AmqpSubscriber{}
	amqp_subscriber.Init(amqp_uri, "elroy.analytics.v0", analytics_client.handle_analytics)	
}
