package controller

func TicketCreateHandler1(subject string, desc string) {
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
	if success := network.TicketAPI(ticket); success {
		//cool it down
		specificsync.CoolItDown(ticket.Id, ticket.Updated)
	}
}
