package model

import (
	"encoding/json"
	"github.com/sankarvj/sample_syncadapter_client/goclient/api"
	"github.com/sankarvj/sample_syncadapter_client/goclient/network"
	"github.com/sankarvj/syncadapter/core"
	"github.com/sankarvj/syncadapter/performer"
	"log"
	"time"
)

type Note struct {
	Ticketid int64 `rt:"tickets" rk:"id"`
	Name     string
	Desc     string `json:"Description"`
	created  time.Time
	core.BaseModel
}

func (note *Note) Create(callback ParallelClientCallback) {
	pro := performer.CreatePro(InitDB())
	//store it local
	note.Id = StoreNote(note)
	log.Println("Local scope -- ", note)

	//localnote := *note //save local instance before passing it to cook for server

	channel := pro.ApiMeltDown(note)
	log.Println("Server scope -- ", note)
	response := createNoteAPI(note)
	if response.Id == network.ResponseSuccess {
		note, err := ParseNote(response.Outcome[0])
		if err != nil {
			close(channel)
			return
		}
		log.Println("channel received")
		channel <- &note
		callback.OnResponseUpdated()
	} else {
		close(channel)
	}
}

func TicketDetail(callback ParallelClientCallback, ticketid int64) {
	pro := performer.CreatePro(InitDB())
	//LOCAL
	dbnotes := ReadNotes(ticketid)
	out, _ := json.Marshal(dbnotes)
	callback.OnResponseReceived(string(out))
	//API
	outcome := api.NotelistAPI(pro.HotId("tickets", ticketid))
	notes, _ := ParseNotes(outcome)
	newnotes, modifiednotes := pro.WhatToDoLogic1(notes, performer.PasserSlice(dbnotes))
	newnotesparsed, _ := ParseNotes(newnotes)
	modifiednotesnotesparsed, _ := ParseNotes(modifiednotes)
	for i := 0; i < len(newnotesparsed); i++ {
		StoreNote(&newnotesparsed[i])
	}

	for i := 0; i < len(modifiednotesnotesparsed); i++ {
		UpdateNote(&modifiednotesnotesparsed[i])
	}

	if pro.DatabaseChanged {
		callback.OnResponseUpdated()
	}
}

//Signal
func (note *Note) Signal(technique core.Technique) bool {
	switch technique {
	case core.BASIC_CREATE:

	case core.BASIC_UPDATE:

	}
	return false
}

func createNoteAPI(note *Note) network.Response {
	var errorresponse = new(network.Response)
	jsonbody, err := json.Marshal(note)
	if err != nil {
		log.Println("Can't marshal note")
		return *errorresponse
	}
	outcome, success := api.CreatNoteAPI(jsonbody)
	if !success {
		return *errorresponse
	}
	*note, err = ParseNote(outcome)
	if err != nil {
		log.Println("Error parsing note")
		return *errorresponse
	}
	errorresponse.Id = network.ResponseSuccess
	errorresponse.Outcome = append(errorresponse.Outcome, *note)
	return *errorresponse
}
