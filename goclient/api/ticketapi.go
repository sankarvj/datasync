package api

import (
	"github.com/sankarvj/sample_syncadapter_client/goclient/network"
)

func CreateTicketAPI(jsonbody []byte) (interface{}, bool) {
	done := make(chan bool)
	var response network.Response
	go func() {
		response = network.MakeCallToServer(network.METHOD_POST, "tickets", jsonbody)
		done <- true
	}()
	<-done
	if len(response.Outcome) > 0 {
		return response.Outcome[0], true
	} else {
		return nil, false
	}
}

func EditTicketAPI(jsonbody []byte) (interface{}, bool) {
	done := make(chan bool)
	var response network.Response
	go func() {
		response = network.MakeCallToServer(network.METHOD_PUT, "tickets", jsonbody)
		done <- true
	}()
	<-done
	if len(response.Outcome) > 0 {
		return response.Outcome[0], true
	} else {
		return nil, false
	}
}

func TicketlistAPI() interface{} {
	done := make(chan bool)
	var response network.Response
	go func() {
		response = network.MakeCallToServer(network.METHOD_GET, "tickets", nil)
		done <- true
	}()
	<-done
	if len(response.Outcome) > 0 {
		return response.Outcome[0]
	} else {
		return nil
	}
}
