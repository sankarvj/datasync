package model

import (
	"encoding/json"
)

//Ticket
func ParseTicket(response interface{}) (*Ticket, error) {
	var obj *Ticket
	out, err := json.Marshal(response)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(out, &obj)
	return obj, err
}
