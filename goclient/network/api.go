package network

import (
	"encoding/json"
	"gitlab.com/vjopensrc/datasync/goclient/model"
	"log"
)

func Sample(ticket *model.Ticket) {
	jsonbody, err := json.Marshal(ticket)
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

	ticket, err = model.ParseTicket(response.Outcome[0])

	if err != nil {
		log.Println("Error parsing ticket")
	}

	log.Println("ticket ", ticket)
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
