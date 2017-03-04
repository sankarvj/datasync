package controller

import (
	"gitlab.com/vjopensrc/datasync/goclient/adapter"
	"gitlab.com/vjopensrc/datasync/goclient/model"
	"gitlab.com/vjopensrc/datasync/goclient/network"
	"log"
)

func TicketCreateHandler(subject string, desc string) {
	//create a ticket object
	ticket := new(model.Ticket)
	ticket.Subject = subject
	ticket.Desc = desc

	specificsync := adapter.CreateSpecificSyncer(model.InitDB())
	//store it local
	specificsync.MakeLocal(model.StoreTicket, ticket)
	//server annex
	//localticket := *ticket //save local instance before passing it to cook for server
	specificsync.CookForRemote(ticket)
	//call api while the object it hot
	network.TicketAPI(ticket)
	//cool it down
	specificsync.CoolItDown(ticket.Id, ticket.Updated)
}

func NoteCreateHandler(name string, desc string, ticketid int64) {
	//create a note object
	note := new(model.Note)
	note.Ticketid = ticketid
	note.Name = name
	note.Desc = desc

	specificsync := adapter.CreateSpecificSyncer(model.InitDB())
	//store it local
	specificsync.MakeLocal(model.StoreNote, note)
	//server call
	//localnote := *note //save local instance before passing it to cook for server
	log.Println("local id ", note.Id)
	log.Println("local ticketid ", note.Ticketid)
	specificsync.CookForRemote(note)
	log.Println("server id ", note.Id)
	log.Println("server ticketid ", note.Ticketid)
	//call api while it is hot
	network.NoteAPI(note)
	//update local with key,sync and updatedtime
	specificsync.CoolItDown(note.Id, note.Updated)
}

func TicketListHandler() {
	specificsync := adapter.CreateSpecificSyncer(model.InitDB())
	//LOCAL
	dbtickets := model.ReadTickets()
	//API
	tickets := network.TicketlistAPI()
	for i := 0; i < len(tickets); i++ {
		ticket := &tickets[i]
		//HOT to COLD conversion
		dbid, dowhat := specificsync.WhatToDo(ticket, adapter.PasserSlice(dbtickets))
		switch dowhat {
		case adapter.CREATE:
			specificsync.MakeLocal(model.StoreTicket, ticket)
			break
		case adapter.UPDATE:
			model.UpdateTicket(ticket, dbid)
			break
		}
	}
}

func NoteListHandler(ticketid int64) {
	specificsync := adapter.CreateSpecificSyncer(model.InitDB())
	//LOCAL
	dbnotes := model.ReadNotes(ticketid)
	//API
	notes := network.NotelistAPI(specificsync.HotId("tickets", ticketid))
	for i := 0; i < len(notes); i++ {
		note := &notes[i]
		//HOT to COLD conversion
		dbid, dowhat := specificsync.WhatToDo(note, adapter.PasserSlice(dbnotes))
		switch dowhat {
		case adapter.CREATE:
			specificsync.MakeLocal(model.StoreNote, note)
			break
		case adapter.UPDATE:
			model.UpdateNote(note, dbid)
			break
		}
	}
}
