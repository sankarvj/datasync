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

	specificsync := adapter.Specificsync{model.InitDB(), "", 0, 0}
	//store it local
	specificsync.CreateLocal(model.StoreTicket, ticket)
	//server annex
	specificsync.AnnexRemote(ticket)
	log.Println("serverticket ", ticket)
	//call api
	network.Sample(ticket)
}

func NoteCreateHandler(name string, desc string, ticketid int64) {
	//create a note object
	note := new(model.Note)
	note.Ticketid = ticketid
	note.Name = name
	note.Desc = desc

	specificsync := adapter.Specificsync{model.InitDB(), "", 0, 0}
	//store it local
	specificsync.CreateLocal(model.StoreNote, note)
	//server call

	log.Println("A localnote ", note)
	localnote := *note
	specificsync.AnnexRemote(note)
	log.Println("B localnote ", note)
	log.Println("C localnote ", localnote)

}
