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
	log.Println("local tickets ::: ", dbtickets)
	//API
	tickets := network.TicketlistAPI()
	log.Println("server tickets ::: ", tickets)
	//HOT to COLD conversion
	for i := 0; i < len(tickets); i++ {
		ticket := &tickets[i]
		log.Println("hot ticket ::: ", ticket)
		specificsync.CookFromRemote(ticket)
		//Not fully cold still it is little hot. The id has the reference to the serverkey
		log.Println("cold ticket ::: ", ticket)

		index := specificsync.FindLocalItemIndex(ticket.Id, adapter.PasserSlice(dbtickets))
		log.Println("index :: ", index)
		if index != -1 { //already stored
			if adapter.NeedUpdate(ticket.Updated, dbtickets[index].Updated) {
				//your update logic
				log.Println("you need to update this :: ", index)
			}
		} else { // new entry
			specificsync.MakeLocal(model.StoreTicket, ticket)
		}

	}

}
