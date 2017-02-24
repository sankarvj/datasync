package goclient

func TicketCreateHandler(subject string, desc string) {

	//initialize database
	db := initDB()
	//create a ticket
	ticket := new(Ticket)
	ticket.subject = subject
	ticket.desc = desc
	//store it locally
	storeTicket(db, ticket)

}
