package model

import (
	"encoding/json"
	"gitlab.com/vjopensrc/datasync/goclient/api"
	"gitlab.com/vjsideprojects/seekite_client_logic/syncadapter/core"
	"gitlab.com/vjsideprojects/seekite_client_logic/syncadapter/performer"
	"log"
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

func (ticket *Ticket) Create(callback ParallelClientCallback) {
	pro := performer.CreatePro(InitDB())
	//store it local
	pro.Prepare(StoreTicket, ticket)
	out, _ := json.Marshal(ticket)
	callback.OnResponseReceived(string(out))
	//localticket := *ticket //save local instance before passing it to cook for server
	if success := pro.Push(ticket); success {
		callback.OnResponseUpdated()
	}
}

func (ticket *Ticket) Update(callback ParallelClientCallback) {
	pro := performer.CreatePro(InitDB())
	//store it local
	pro.Prepare(UpdateTicket, ticket)
	out, _ := json.Marshal(ticket)
	callback.OnResponseReceived(string(out))
	//localticket := *ticket //save local instance before passing it to cook for server
	if success := pro.Push(ticket); success {
		callback.OnResponseUpdated()
	}
}

func TicketList(callback ParallelClientCallback) {
	pro := performer.CreatePro(InitDB())
	//LOCAL
	dbtickets := ReadTickets()
	out, _ := json.Marshal(dbtickets)
	callback.OnResponseReceived(string(out))
	//API
	outcome := api.TicketlistAPI()
	tickets, _ := ParseTickets(outcome)

	newitems, updateditems := pro.WhatToDoLogic1(tickets, performer.PasserSlice(dbtickets))

	newTickets, _ := ParseTickets(newitems)
	updatedTickets, _ := ParseTickets(updateditems)
	storeTickets(newTickets)
	modifyTickets(updatedTickets)
	//store this to db and send call
	if pro.DatabaseChanged {
		callback.OnResponseUpdated()
	}

	//TODO MOVE THIS INSIDE ADAPTER
	// for i := 0; i < len(tickets); i++ {
	// 	ticket := &tickets[i]
	// 	//HOT to COLD conversion
	// 	dowhat := pro.WhatToDo(ticket, performer.PasserSlice(dbtickets))
	// 	switch dowhat {
	// 	case performer.CREATE:
	// 		databasechanged = true
	// 		pro.Prepare(StoreTicket, ticket)
	// 		break
	// 	case performer.UPDATE:
	// 		databasechanged = true
	// 		UpdateTicket(ticket)
	// 		break
	// 	}
	// }

}

func storeTickets(newtickets []Ticket) {
	for i := 0; i < len(newtickets); i++ {
		newticket := newtickets[i]
		log.Println("newticket ::: ", newticket.Key)
		StoreTicket(&newticket)
	}
}

func modifyTickets(modifiedtickets []Ticket) {
	for i := 0; i < len(modifiedtickets); i++ {
		modifiedticket := modifiedtickets[i]
		log.Println("modifiedticket ::: ", modifiedticket.Id)
		UpdateTicket(&modifiedticket)
	}
}

func (ticket *Ticket) Signal(technique int) bool {
	var success bool
	switch technique {
	case performer.TECHNIQUE_BASIC_CREATE:
		success = createAPI(ticket)
		break
	case performer.TECHNIQUE_BASIC_UPDATE:
		success = updateAPI(ticket)
		break
	case performer.TECHNIQUE_PERIODIC_SHOT:
		syncFrozenData()
		break
	}
	return success
}

func createAPI(ticket *Ticket) bool {
	jsonbody, err := json.Marshal(ticket)
	if err != nil {
		log.Println("Can't marshal ticket")
		return false
	}
	outcome, success := api.CreateTicketAPI(jsonbody)
	if !success {
		return false
	}
	*ticket, err = ParseTicket(outcome)
	if err != nil {
		log.Println("Error parsing ticket")
		return false
	}
	return true
}

func updateAPI(ticket *Ticket) bool {
	jsonbody, err := json.Marshal(ticket)
	if err != nil {
		log.Println("Can't marshal ticket")
		return false
	}
	outcome, success := api.EditTicketAPI(jsonbody)
	if !success {
		return false
	}
	*ticket, err = ParseTicket(outcome)
	if err != nil {
		log.Println("Error parsing ticket")
		return false
	}
	return true
}

func syncFrozenData() {
	dbtickets := ReadFrozenTickets()
	for i := 0; i < len(dbtickets); i++ {
		pro := performer.CreatePro(InitDB())
		ticket := &dbtickets[i]
		pro.Push(ticket)
	}
}
