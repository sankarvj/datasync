package model

import (
	"gitlab.com/vjopensrc/datasync/goclient/adapter"
	"time"
)

type Ticket struct {
	Subject   string
	Desc      string `rt:"trips" rk:"id"`
	requester string
	agent     string
	created   time.Time
	adapter.BaseModel
}

type Note struct {
	Ticketid int64 `rt:"tickets" rk:"id"`
	Name     string
	Desc     string
	created  time.Time
	adapter.BaseModel
}
