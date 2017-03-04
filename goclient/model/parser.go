package model

import (
	"encoding/json"
)

//Ticket
func ParseTicket(response interface{}) (Ticket, error) {
	var obj *Ticket
	out, err := json.Marshal(response)
	if err != nil {
		return *obj, err
	}
	err = json.Unmarshal(out, &obj)
	return *obj, err
}

//Note
func ParseNote(response interface{}) (Note, error) {
	var obj *Note
	out, err := json.Marshal(response)
	if err != nil {
		return *obj, err
	}
	err = json.Unmarshal(out, &obj)
	return *obj, err
}

//Tickets
func ParseTickets(response interface{}) ([]Ticket, error) {
	tickets := make([]Ticket, 0)
	out, err := json.Marshal(response)
	if err != nil {
		return tickets, err
	}
	err = json.Unmarshal(out, &tickets)
	return tickets, err
}

//Notes
func ParseNotes(response interface{}) ([]Note, error) {
	notes := make([]Note, 0)
	out, err := json.Marshal(response)
	if err != nil {
		return notes, err
	}
	err = json.Unmarshal(out, &notes)
	return notes, err
}
