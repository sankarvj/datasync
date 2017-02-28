package controller

import (
	"gitlab.com/vjopensrc/datasync/goclient/adapter"
	"gitlab.com/vjopensrc/datasync/goclient/model"
)

func TicketCreateHandler(subject string, desc string) {
	//create a ticket object
	ticket := new(model.Ticket)
	ticket.Subject = subject
	ticket.Desc = desc

	//store it local
	adapter.CreateLocal(model.StoreTicket, ticket)

	//server call
	db := model.InitDB()
	adapter.ProcessForRemote(db, ticket)
}
