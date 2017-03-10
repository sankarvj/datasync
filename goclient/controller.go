package goclient

import (
	"encoding/json"
	"gitlab.com/vjopensrc/datasync/syncadapter/performer"
	"gitlab.com/vjopensrc/datasync/syncadapter/technique"
	"log"
)

func TicketCreateHandler(callback ClientCallback, subject string, desc string) {
	//create a ticket object
	ticket := new(Ticket)
	ticket.Subject = subject
	ticket.Desc = desc

	pro := performer.CreatePro(InitDB())
	//store it local
	pro.Prepare(StoreTicket, ticket)
	out, _ := json.Marshal(ticket)
	callback.OnResponseReceived(string(out))
	//server annex
	//localticket := *ticket //save local instance before passing it to cook for server
	pro.CookForRemote(ticket)
	//call api while the object it hot
	if success := pro.CallRemote(ticket); success {
		callback.OnResponseUpdated()
	}
}

func TicketEditHandler(subject string, desc string, ticketid int64) {
	//create a ticket object
	ticket := new(Ticket)
	ticket.Id = ticketid
	ticket.Subject = subject
	ticket.Desc = desc

	pro := performer.CreatePro(InitDB())
	//store it local
	pro.Prepare(UpdateTicket, ticket)
	//server annex
	//localticket := *ticket //save local instance before passing it to cook for server
	log.Println("local id ", ticket.Id)
	pro.CookForRemote(ticket)
	log.Println("server id ", ticket.Id)

	if success := pro.CallRemote(ticket); success {

	}

}

func TicketListHandler(callback ClientCallback) {
	pro := performer.CreatePro(InitDB())
	//LOCAL
	dbtickets := ReadTickets()
	out, _ := json.Marshal(dbtickets)
	callback.OnResponseReceived(string(out))
	//API
	databasechanged := false
	tickets := TicketlistAPI()
	for i := 0; i < len(tickets); i++ {
		ticket := &tickets[i]
		//HOT to COLD conversion
		dowhat := pro.WhatToDo(ticket, performer.PasserSlice(dbtickets))
		switch dowhat {
		case performer.CREATE:
			databasechanged = true
			pro.Prepare(StoreTicket, ticket)
			break
		case performer.UPDATE:
			databasechanged = true
			UpdateTicket(ticket)
			break
		}
	}

	if databasechanged {
		callback.OnResponseUpdated()
	}
}

func NoteCreateHandler(name string, desc string, ticketid int64) {
	//create a note object
	note := new(Note)
	note.Ticketid = ticketid
	note.Name = name
	note.Desc = desc

	pro := performer.CreatePro(InitDB())
	//store it local
	pro.Prepare(StoreNote, note)
	//server call
	//localnote := *note //save local instance before passing it to cook for server
	log.Println("local id ", note.Id)
	log.Println("local ticketid ", note.Ticketid)
	pro.CookForRemote(note)
	log.Println("server id ", note.Id)
	log.Println("server ticketid ", note.Ticketid)
	//call api while it is hot
	if success := pro.CallRemote(note); success {

	}
}

func NoteListHandler(callback ClientCallback, ticketid int64) {
	pro := performer.CreatePro(InitDB())
	//LOCAL
	dbnotes := ReadNotes(ticketid)
	out, _ := json.Marshal(dbnotes)
	callback.OnResponseReceived(string(out))
	//API
	databasechanged := false
	notes := NotelistAPI(pro.HotId("tickets", ticketid))
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

func PeriodicSync() {
	periodicsync := technique.CreatePeriodic(InitDB())
	periodicsync.Models = append(periodicsync.Models, &Ticket{})
	periodicsync.Models = append(periodicsync.Models, &Note{})
	periodicsync.CheckPeriodic()
}
