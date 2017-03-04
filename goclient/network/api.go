package network

import (
	"encoding/json"
	"gitlab.com/vjopensrc/datasync/goclient/model"
	"log"
	"strconv"
)

func TicketAPI(ticket *model.Ticket) {
	jsonbody, err := json.Marshal(ticket)
	log.Println("jsonbody ::: ", jsonbody)
	if err != nil {
		log.Println("Can't marshal ticket")
		return
	}
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_post, "tickets", jsonbody)
		done <- true
	}()
	<-done

	*ticket, err = model.ParseTicket(response.Outcome[0])

	if err != nil {
		log.Println("Error parsing ticket")
	}

	log.Println("ticket ", ticket)
}

func NoteAPI(note *model.Note) {
	jsonbody, err := json.Marshal(note)
	if err != nil {
		log.Println("Can't marshal note")
		return
	}
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_post, "notes", jsonbody)
		done <- true
	}()
	<-done

	*note, err = model.ParseNote(response.Outcome[0])

	if err != nil {
		log.Println("Error parsing note")
	}
}

func TicketlistAPI() []model.Ticket {
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_get, "tickets", nil)
		done <- true
	}()
	<-done
	tickets, _ := model.ParseTickets(response.Outcome[0])
	return tickets
}

func NotelistAPI(ticketid int64) []model.Note {
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_get, "tickets/"+strconv.FormatInt(ticketid, 10), nil)
		done <- true
	}()
	<-done
	notes, _ := model.ParseNotes(response.Outcome[1])
	return notes
}

// func createTicketAPI(ticket *model.Ticket) {
// 	//Prepare the object
// 	db = initDB()
// 	oldticket := *ticket
// 	jsonbody, err := shapeMaker(db, ticket, &oldticket, servertripid, lastid)
// 	if err != nil {
// 		callback.OnError(ResponseInternalError, "Error shaping ticket")
// 		return
// 	}
// 	//Make API call - !task object is loaded with server ids
// 	response := makeCallToServer(db, method_post, "tickets", jsonbody)
// 	//Parse response from server
// 	trip, err = parseTrip(response.Outcome[0])
// 	if err != nil {
// 		log.Println("error in ticket create api ", err)
// 		return
// 	}
// 	//Send success callback if found
// 	solved := shapeSolver(db, response, trip, &oldtrip, "trips", callback)
// 	//task object is loaded with db ids
// 	if solved {
// 		//TODO update_trip_logic
// 		callback.OnResponseReceived(structToStr(trip))
// 	}
// }
