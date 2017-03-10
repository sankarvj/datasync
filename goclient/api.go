package goclient

import (
	"encoding/json"
	"log"
	"strconv"
)

func ticketAPI(ticket *Ticket) bool {
	jsonbody, err := json.Marshal(ticket)
	if err != nil {
		log.Println("Can't marshal ticket")
		return false
	}
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_post, "tickets", jsonbody)
		done <- true
	}()
	<-done

	if len(response.Outcome) > 0 {
		*ticket, err = ParseTicket(response.Outcome[0])
		if err != nil {
			log.Println("Error parsing ticket")
			return false
		}
		log.Println("ticket ", ticket)
		return true
	} else {
		return false
	}

}

func ticketEditAPI(ticket *Ticket) bool {
	jsonbody, err := json.Marshal(ticket)
	if err != nil {
		log.Println("Can't marshal ticket")
		return false
	}
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_put, "tickets", jsonbody)
		done <- true
	}()
	<-done

	if len(response.Outcome) > 0 {
		*ticket, err = ParseTicket(response.Outcome[0])
		if err != nil {
			log.Println("Error parsing ticket")
			return false
		}

		log.Println("ticket ", ticket)
		return true
	} else {
		return false
	}

}

func NoteAPI(note *Note) bool {
	jsonbody, err := json.Marshal(note)
	if err != nil {
		log.Println("Can't marshal note")
		return false
	}
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_post, "notes", jsonbody)
		done <- true
	}()
	<-done

	if len(response.Outcome) > 0 {
		*note, err = ParseNote(response.Outcome[0])
		if err != nil {
			log.Println("Error parsing note")
			return false
		}
		return true
	} else {
		return false
	}

}

func TicketlistAPI() []Ticket {
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_get, "tickets", nil)
		done <- true
	}()
	<-done
	if len(response.Outcome) > 0 {
		tickets, _ := ParseTickets(response.Outcome[0])
		return tickets
	} else {
		tickets := make([]Ticket, 0)
		return tickets
	}

}

func NotelistAPI(ticketid int64) []Note {
	done := make(chan bool)
	var response Response
	go func() {
		response = makeCallToServer(method_get, "tickets/"+strconv.FormatInt(ticketid, 10), nil)
		done <- true
	}()
	<-done
	if len(response.Outcome) > 1 {
		notes, _ := ParseNotes(response.Outcome[1])
		return notes
	} else {
		notes := make([]Note, 0)
		return notes
	}

}
