package goclient

import (
	"gitlab.com/vjopensrc/datasync/adapter"
	"time"
)

type Ticket struct {
	Subject   string
	Desc      string `json:"Description"`
	requester string
	agent     string
	created   time.Time
	adapter.BaseModel
}

type Note struct {
	Ticketid int64 `rt:"tickets" rk:"id"`
	Name     string
	Desc     string `json:"Description"`
	created  time.Time
	adapter.BaseModel
}
