package goclient

import (
	"gitlab.com/vjopensrc/datasync/syncadapter/core"
	"gitlab.com/vjopensrc/datasync/syncadapter/performer"
	"time"
)

type Ticket struct {
	Subject   string
	Desc      string `json:"Description"`
	requester string
	agent     string
	created   time.Time
	core.BaseModel
}

func (ticket *Ticket) Signal(technique int) bool {
	var success bool
	switch technique {
	case performer.TECHNIQUE_BASIC_CREATE:
		success = ticketAPI(ticket)
		break
	case performer.TECHNIQUE_BASIC_UPDATE:
		success = ticketEditAPI(ticket)
		break
	case performer.TECHNIQUE_PERIODIC_SHOT:
		syncFrozenData()
		break
	}

	return success
}

func syncFrozenData() {
	dbtickets := ReadFrozenTickets()
	for i := 0; i < len(dbtickets); i++ {
		pro := performer.CreatePro(InitDB())
		ticket := &dbtickets[i]
		pro.CookForRemote(ticket)
		pro.CallRemote(ticket)
	}
}

type Note struct {
	Ticketid int64 `rt:"tickets" rk:"id"`
	Name     string
	Desc     string `json:"Description"`
	created  time.Time
	core.BaseModel
}

func (note *Note) Signal(technique int) bool {
	var success bool
	switch technique {
	case performer.TECHNIQUE_BASIC_CREATE:
		success = NoteAPI(note)
		break
	}

	return success
}
