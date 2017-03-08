package goclient

import (
	"encoding/json"
	"gitlab.com/vjopensrc/datasync/adapter"
	"log"
)

func TicketCreateHandler(callback ClientCallback, subject string, desc string) {
	//create a ticket object
	ticket := new(Ticket)
	ticket.Subject = subject
	ticket.Desc = desc

	specificsync := adapter.CreateSpecificSyncer(InitDB())
	//store it local
	specificsync.MakeLocal(StoreTicket, ticket)
	out, _ := json.Marshal(ticket)
	callback.OnResponseReceived(string(out))
	//server annex
	//localticket := *ticket //save local instance before passing it to cook for server
	specificsync.CookForRemote(ticket)
	//call api while the object it hot
	if success := ticketAPI(ticket); success {
		//cool it down
		specificsync.CoolItDown(ticket.Id, ticket.Updated)
		callback.OnResponseUpdated()
	}
}

func TicketEditHandler(subject string, desc string, ticketid int64) {
	//create a ticket object
	ticket := new(Ticket)
	ticket.Id = ticketid
	ticket.Subject = subject
	ticket.Desc = desc

	specificsync := adapter.CreateSpecificSyncer(InitDB())
	//store it local
	specificsync.UpdateLocal(UpdateTicket, ticket, ticketid)
	//server annex
	//localticket := *ticket //save local instance before passing it to cook for server
	log.Println("local id ", ticket.Id)
	specificsync.CookForRemote(ticket)
	log.Println("server id ", ticket.Id)

	//call api while the object is hot
	if success := ticketEditAPI(ticket); success {
		//cool it down
		specificsync.CoolItDown(ticket.Id, ticket.Updated)
	}

}

func TicketListHandler(callback ClientCallback) {
	specificsync := adapter.CreateSpecificSyncer(InitDB())
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
		dbid, dowhat := specificsync.WhatToDo(ticket, adapter.PasserSlice(dbtickets))
		switch dowhat {
		case adapter.CREATE:
			databasechanged = true
			specificsync.MakeLocal(StoreTicket, ticket)
			break
		case adapter.UPDATE:
			databasechanged = true
			UpdateTicket(ticket, dbid)
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

	specificsync := adapter.CreateSpecificSyncer(InitDB())
	//store it local
	specificsync.MakeLocal(StoreNote, note)
	//server call
	//localnote := *note //save local instance before passing it to cook for server
	log.Println("local id ", note.Id)
	log.Println("local ticketid ", note.Ticketid)
	specificsync.CookForRemote(note)
	log.Println("server id ", note.Id)
	log.Println("server ticketid ", note.Ticketid)
	//call api while it is hot
	if success := NoteAPI(note); success {
		//update local with key,sync and updatedtime
		specificsync.CoolItDown(note.Id, note.Updated)
	}

}

func NoteListHandler(callback ClientCallback, ticketid int64) {
	specificsync := adapter.CreateSpecificSyncer(InitDB())
	//LOCAL
	dbnotes := ReadNotes(ticketid)
	out, _ := json.Marshal(dbnotes)
	callback.OnResponseReceived(string(out))
	//API
	databasechanged := false
	notes := NotelistAPI(specificsync.HotId("tickets", ticketid))
	for i := 0; i < len(notes); i++ {
		note := &notes[i]
		//HOT to COLD conversion
		dbid, dowhat := specificsync.WhatToDo(note, adapter.PasserSlice(dbnotes))
		switch dowhat {
		case adapter.CREATE:
			databasechanged = true
			specificsync.MakeLocal(StoreNote, note)
			break
		case adapter.UPDATE:
			databasechanged = true
			UpdateNote(note, dbid)
			break
		}
	}

	if databasechanged {
		callback.OnResponseUpdated()
	}
}

func GenericSync() {
	genericsync := adapter.CreateGenericSyncer(InitDB())
	genericsync.Tablenames = append(genericsync.Tablenames, "tickets")
	genericsync.SyncFrozenData()
}
