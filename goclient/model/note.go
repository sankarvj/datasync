package model

import (
	"encoding/json"
	"gitlab.com/vjopensrc/datasync/goclient/api"
	"gitlab.com/vjopensrc/datasync/syncadapter/core"
	"gitlab.com/vjopensrc/datasync/syncadapter/performer"
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
	pro.Prepare(StoreNote, note)
	//localnote := *note //save local instance before passing it to cook for server
	if success := pro.Push(note); success {
		callback.OnResponseUpdated()
	}
}

func TicketDetail(callback ParallelClientCallback, ticketid int64) {
	pro := performer.CreatePro(InitDB())
	//LOCAL
	dbnotes := ReadNotes(ticketid)
	out, _ := json.Marshal(dbnotes)
	callback.OnResponseReceived(string(out))
	//API
	databasechanged := false
	outcome := api.NotelistAPI(pro.HotId("tickets", ticketid))
	notes, _ := ParseNotes(outcome)
	for i := 0; i < len(notes); i++ {
		note := &notes[i]
		//HOT to COLD conversion
		dowhat := pro.WhatToDo(note, performer.PasserSlice(dbnotes))
		switch dowhat {
		case performer.CREATE:
			databasechanged = true
			pro.Prepare(StoreNote, note)
			break
		case performer.UPDATE:
			databasechanged = true
			UpdateNote(note)
			break
		}
	}

	if databasechanged {
		callback.OnResponseUpdated()
	}
}

func (note *Note) Signal(technique int) bool {
	var success bool
	switch technique {
	case performer.TECHNIQUE_BASIC_CREATE:
		success = createNoteAPI(note)
		break
	}

	return success
}

func createNoteAPI(note *Note) bool {
	jsonbody, err := json.Marshal(note)
	if err != nil {
		log.Println("Can't marshal note")
		return false
	}
	outcome, success := api.CreatNoteAPI(jsonbody)
	if !success {
		return false
	}
	*note, err = ParseNote(outcome)
	if err != nil {
		log.Println("Error parsing note")
		return false
	}
	return true
}
