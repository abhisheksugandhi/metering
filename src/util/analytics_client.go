package util

import(
    "fmt"
    "net/http"
    "bytes"
    "io"
    "io/ioutil"
)

type AnalyticsClient struct{	
    url string;
    client http.Client
}

func (a *AnalyticsClient) Init(serverUrl string) {
    a.url = serverUrl
    a.client = http.Client{}
}

func (a *AnalyticsClient) Call(data []byte) {
    req, err := http.NewRequest("POST", a.url, bytes.NewBuffer(data))
    req.Header.Set("Content-Type", "application/octet-stream")
    resp, err := a.client.Do(req)
    if err != nil {
	fmt.Errorf("Server call error %s", err) 
    }
    defer resp.Body.Close()
    io.Copy(ioutil.Discard, res.Body)
}
