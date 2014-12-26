package main

import (
	"bytes"
	"encoding/json"
	"flywheel/util"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
)

type GDMClient struct {
	url       string
	client    http.Client
	publisher util.AmqpPublisher
	subsciber util.AmqpSubscriber
}

type GDMClientRequest struct {
	Id           int64
	Origin       string
	Destinations string
}

type Location struct {
	Lat float64
	Lon float64
}

type GDMResponse struct {
	Status string
	Rows   []GDMResponseRow
}

type GDMResponseRow struct {
	Elements []GDMResponseElement
}

type GDMResponseElement struct {
	Status   string
	Duration GDMResponseDurationElement
	Distance GDMResponseDistanceElement
}

type GDMResponseDurationElement struct {
	Value int64
	Text  string
}

type GDMResponseDistanceElement struct {
	Value int64
	Text  string
}

type ResponseElement struct {
	Duration int64
	Distance int64
}

type GDMClientResponse struct {
	Id       int64 `json:"id"`
	Status   string
	Response []ResponseElement
}

func (a *GDMClient) Init(serverUrl string, amqpUri string) {
	a.url = serverUrl
	a.client = http.Client{}
	a.publisher.Init(amqpUri)
}

func (a *GDMClient) subscribe_to_queue(amqpUri string, queueName string) {
	a.subsciber = util.AmqpSubscriber{}
	a.subsciber.Init(amqpUri, queueName, a.handle_eta)
}

func (a *GDMClient) GetDM(data []byte, replyTo string, messageId string, correlationId string) {
	var inreq GDMClientRequest
	err := json.Unmarshal(data, &inreq)

	if err != nil {
		panic(err)
	}

	var queryParams bytes.Buffer
	queryParams.WriteString(fmt.Sprintf("origins=%s", inreq.Origin))
	queryParams.WriteString(fmt.Sprintf("&destinations=%s", inreq.Destinations))

	fmt.Println("calling server url " + fmt.Sprintf("%s?%s", a.url, queryParams.String()))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", a.url, queryParams.String()), nil)
	resp, err := a.client.Do(req)
	if err != nil {
		fmt.Errorf("Server call error %s", err)
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)

	var apiResponse GDMResponse
	json.Unmarshal(responseBody, &apiResponse)
	fmt.Printf("got response from GDM server as %s \n", bytes.NewBuffer(responseBody).String())

	var clientResponse GDMClientResponse
	if apiResponse.Status == "OK" {
		var responseElementArr = make([]ResponseElement, 0)
		elements := apiResponse.Rows[0].Elements
		for _, element := range elements {
			if element.Status == "OK" {
				responseElementArr = append(responseElementArr, ResponseElement{element.Duration.Value, element.Distance.Value})
			} else {
				responseElementArr = append(responseElementArr, ResponseElement{-1, -1})
			}
		}
		clientResponse = GDMClientResponse{inreq.Id, "OK", responseElementArr}
	} else {
		var responseElementArr = make([]ResponseElement, 0)
		clientResponse = GDMClientResponse{inreq.Id, "ERROR", responseElementArr}
	}

	var responseData, _ = json.Marshal(&clientResponse)
	fmt.Printf("Response to sent to %s after calling GDM %s\n", replyTo, bytes.NewBuffer(responseData).String())
	err1 := a.publisher.Publish_Amqp(responseData, replyTo, messageId, correlationId)
	if err != nil {
		fmt.Printf("Exchange Publish: %s", err1)
	}
}

func (a *GDMClient) handle_eta(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		d.Ack(false)
		log.Printf("got event of length %d: %q", len(d.Body), d.Body)
		a.GetDM(d.Body, d.ReplyTo, d.MessageId, d.CorrelationId)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
