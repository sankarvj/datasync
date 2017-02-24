package goclient

import (
	"time"
)

//basemodel interface
type Basemodel interface {
	//Key
	getKey() int64
	//Id
	getId() int64
	setId(id int64)
	//Sync
	getSynced() bool
	setSynced(sync bool)
	//Time
	getUpdatedat() int64
}

type Ticket struct {
	id        int64
	subject   string
	desc      string
	requester string
	agent     string
	updated   int64
	created   time.Time
	key       int64
	synced    bool
}
