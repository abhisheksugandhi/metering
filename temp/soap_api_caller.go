package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

const SOAP_URL = "https://demo.taxitronic.org/VTS_SERVICE/group_service.asmx"

const SOAP_QUERY_FORMAT = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
   <GET_CAB_LIST_LAST_INFO xmlns="https://www.taxitronic.org/VTS_SERVICE/">
      <UserId>%s</UserId>
      <Password>%s</Password>
    </GET_CAB_LIST_LAST_INFO>
  </soap:Body>
</soap:Envelope`

type CabInfo struct {
	CabNumber string `xml:"CabNumber" json:"CabNumber"`
	DriverID  string `xml:"DriverID" json:"driverId"`
	FleetCode string `xml:"FleetCode" json:"fleetCode"`
	Latitude  string `xml:"GPS_LA" json:"lat"`
	Longitude string `xml:"GPS_LO" json:"lng"`
}

type SoapBody struct {
	CabListLastInfo []CabInfo `xml:GET_CAB_LIST_LAST_INFOResponse><GET_CAB_LIST_LAST_INFOResult>List`
}

type SoapEnvelope struct {
	XMLName xml.Name
	Body    SoapBody
}

func call_soap_api(config VerifoneAPIConfig) []CabInfo {
	soapRequestContent := fmt.Sprintf(SOAP_QUERY_FORMAT, config.UserName, config.Password)
	httpClient := new(http.Client)
	resp, err := httpClient.Post(SOAP_URL, "text/xml; charset=utf-8", bytes.NewBufferString(soapRequestContent))
	if err != nil {
		// handle error
	}
	b, e := ioutil.ReadAll(resp.Body) // probably not efficient, done because the stream isn't always a pure XML stream and I have to fix things (not shown here)
	if e != nil {
		panic(e)
	}
	in := string(b)
	parser := xml.NewDecoder(bytes.NewBufferString(in))
	envelope := new(SoapEnvelope) // this allocates the structure in which we'll decode the XML
	err = parser.DecodeElement(&envelope, nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	return envelope.Body.CabListLastInfo
}

func init_soap_api_caller() {

}
