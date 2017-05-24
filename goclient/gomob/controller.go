package gomob

import (
	"github.com/sankarvj/sample_syncadapter_client/goclient/model"
	"github.com/sankarvj/syncadapter/technique"
)

func init() {

}

func TicketCreateHandler(callback ClientCallback, subject string, desc string) {
	//create a ticket object
	ticket := new(model.Ticket)
	ticket.Subject = subject
	ticket.Desc = desc
	ticket.Create(callback)
}

func TicketEditHandler(callback ClientCallback, subject string, desc string, ticketid int64) {
	//create a ticket object
	ticket := new(model.Ticket)
	ticket.Id = ticketid
	ticket.Subject = subject
	ticket.Desc = desc
	ticket.Update(callback)
}

func TicketListHandler(callback ClientCallback) {
	model.TicketList(callback)
}

func NoteCreateHandler(callback ClientCallback, name string, desc string, ticketid int64) {
	//create a note object
	note := new(model.Note)
	note.Ticketid = ticketid
	note.Name = name
	note.Desc = desc
	note.Create(callback)
}

func NoteListHandler(callback ClientCallback, ticketid int64) {
	model.TicketDetail(callback, ticketid)
}

func PeriodicSync() {
	periodicsync := technique.CreatePeriodic(model.InitDB())
	periodicsync.Models = append(periodicsync.Models, &model.Ticket{})
	periodicsync.Models = append(periodicsync.Models, &model.Note{})
	periodicsync.CheckPeriodic()
}

func RemoteSync() {
	periodicsync := technique.CreatePeriodic(model.InitDB())
	periodicsync.Models = append(periodicsync.Models, &model.Ticket{})
	periodicsync.Models = append(periodicsync.Models, &model.Note{})
	periodicsync.CheckPeriodic()
}
