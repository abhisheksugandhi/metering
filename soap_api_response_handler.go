package main

import (
	"encoding/json"
)

type Vehicle struct {
	FleetId   int    `json:"fleet_id"`
	CabNumber string `json:"cabNumber"`
}

type MeterOn struct {
	Msgtype string   `json:"type"`
	Subtype string   `json:"subtype"`
	Id      string   `json:"id"`
	Cab     *Vehicle `json:"content"`
}

func init_soap_response_handler() {
	init_fleet_manager()
	init_amqp_publisher()
}

func handle_soap_api_response(cabs_response []CabInfo) {
	mapping := get_fleet_mapping()
	for _, cab := range cabs_response {
		fleet_xid, ok := mapping[cab.FleetCode]
		if ok {
			vehicle := Vehicle{fleet_xid, cab.CabNumber}
			message := MeterOn{"drivers", "meter_on", "", &vehicle}
			data, _ := json.Marshal(&message)
			publish_amqp("spacely.metering.v0", data)
		}
	}
}
