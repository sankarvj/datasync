package api

import (
	"gitlab.com/vjopensrc/datasync/goclient/network"
	"strconv"
)

func CreatNoteAPI(jsonbody []byte) (interface{}, bool) {
	done := make(chan bool)
	var response network.Response
	go func() {
		response = network.MakeCallToServer(network.METHOD_POST, "notes", jsonbody)
		done <- true
	}()
	<-done
	if len(response.Outcome) > 0 {
		return response.Outcome[0], true
	} else {
		return nil, false
	}
}

func NotelistAPI(ticketid int64) interface{} {
	done := make(chan bool)
	var response network.Response
	go func() {
		response = network.MakeCallToServer(network.METHOD_GET, "tickets/"+strconv.FormatInt(ticketid, 10), nil)
		done <- true
	}()
	<-done
	if len(response.Outcome) > 1 {
		return response.Outcome[1]
	} else {
		return nil
	}

}
