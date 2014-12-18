package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AnalyticsClient struct {
	url    string
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
	fmt.Printf("got response status %s on analytics call\n", resp.Status)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("got response from GDM server as %s \n", bytes.NewBuffer(responseBody).String())
}
