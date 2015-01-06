package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type AnalyticsClient struct {
	url        string
	client     http.Client
	apiKey     string
	apiSecret  []byte
	apiVersion string
}

type AnalyticsAPIConfig struct {
	Url     string
	Key     string
	Secret  string
	Version string
}

func (a *AnalyticsClient) Init(serverUrl string) {
	a.url = serverUrl
	a.client = http.Client{}
	a.readAnalyticsAPIConfig()
}

func (a *AnalyticsClient) readAnalyticsAPIConfig() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	data, err := ioutil.ReadFile("config/analytics_api.yml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]AnalyticsAPIConfig)
	err = yaml.Unmarshal(data, &m)
	analyticsApiConfig := m[environment]

	a.url = analyticsApiConfig.Url
	a.apiKey = analyticsApiConfig.Key
	a.apiSecret = []byte(analyticsApiConfig.Secret)
	a.apiVersion = analyticsApiConfig.Version
}

func (a *AnalyticsClient) Call(data []byte) {
	postData := url.Values{}
	content := bytes.NewBuffer(data).String()
	postData.Set("content", content)

	signData := fmt.Sprintf("/add|%s|%s", a.apiVersion, content)
	url := fmt.Sprintf("%s/add?api_ver=%s&api_key=%s&sig=%s", a.url, a.apiVersion, a.apiKey, mac([]byte(signData), a.apiSecret))
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(postData.Encode()))
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

func mac(text, secret []byte) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(text))
	expectedMac := mac.Sum(nil)
	return hex.EncodeToString(expectedMac)
}
