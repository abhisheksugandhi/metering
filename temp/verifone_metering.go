package main

import (
	"fmt"
	"time"
)

func init_main() {
	init_configs()
	init_fleet_manager()
	init_amqp_publisher()
	init_amqp_subscriber()
}

func main() {
	init_main()
	for {
		cabs_info := call_soap_api(get_api_config())
		handle_soap_api_response(cabs_info)
		fmt.Println("sleeping for 1 minute")
		time.Sleep(60 * 1000 * time.Millisecond)
	}
}
