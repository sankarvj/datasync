package performer

import (
	"database/sql"
	"gitlab.com/vjopensrc/datasync/syncadapter/core"
	"reflect"
	"strconv"
)

const (
	NOTHING = 0
	CREATE  = 1
	UPDATE  = 2
)

const (
	TECHNIQUE_BASIC_CREATE  = 0
	TECHNIQUE_BASIC_UPDATE  = 1
	TECHNIQUE_PERIODIC_SHOT = 2
)

type Pro struct {
	DBInst    *sql.DB
	Tablename string
	Localid   int64
}

func CreatePro(db *sql.DB) Pro {
	return Pro{db, "", 0}
}

func CreateProAdv(db *sql.DB, tablename string, localid int64) Pro {
	return Pro{db, tablename, localid}
}

//Expects a function followed by its params
//One of the param must implement cooker interface. This is mandatory.
//Cooker interface is responsible for datasync. Datasync skips if no params implements cooker
func (s *Pro) Prepare(fn interface{}, params ...interface{}) {
	var cooker core.Cooker
	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(params) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		if inImplementsCooker(in) {
			cooker = in.(core.Cooker)
		}
		inputs[k] = reflect.ValueOf(in)
	}

	var performupdate bool
	if cooker != nil {
		if cooker.LocalId() != 0 {
			performupdate = true
		} else {
			//Incase if it created in local then 0 wil be set to Id #noproblem
			//If it is coming from the server then key will be set to Id so key column will be updated
			cooker.SetLocalId(cooker.ServerKey())
		}

	} else {
		oops("Cannot perform datasync no param implements cooker...", true)
		return
	}

	cooker.PrepareLocal(performupdate)

	if performupdate {
		f.Call(inputs)
	} else {
		result := f.Call(inputs)
		//find the last inserted id
		s.Localid = result[0].Interface().(int64)
		cooker.SetLocalId(s.Localid)
	}

}

func (s *Pro) Push(cooker core.Cooker) bool {
	var remoteUpdated bool
	if cooker.ServerKey() != 0 { //update
		remoteUpdated = cooker.Signal(TECHNIQUE_BASIC_UPDATE)
	} else { //create
		remoteUpdated = cooker.Signal(TECHNIQUE_BASIC_CREATE)
	}

	if remoteUpdated {
		//Using LocalId here is very misleading. Also we can't sure that the user implementation always update ID as serverkey
		s.coolItDown(cooker.LocalId(), cooker.UpdatedAt())
	}

	return remoteUpdated
}

func (s *Pro) coolItDown(key int64, updated int64) {
	updateKey(s.DBInst, s.Tablename, key, s.Localid, updated)
}

func (s *Pro) HotId(tablename string, id int64) int64 {
	return serverVal(s.DBInst, tablename, strconv.FormatInt(id, 10))
}

func (s *Pro) WhatToDo(cooker core.Cooker, dblistitems []core.Passer) (dowhat int64) {
	//HOT to COLD conversion
	s.CookFromRemote(cooker)
	index := -1
	dowhat = NOTHING
	for i := 0; i < len(dblistitems); i++ {
		if (dblistitems[i]).ServerKey() == cooker.ServerKey() {
			index = i
		}
	}

	if index != -1 { //already stored
		if needUpdate(cooker.UpdatedAt(), dblistitems[index].UpdatedAt()) {
			dowhat = UPDATE
			cooker.SetLocalId(dblistitems[index].LocalId())
		}
	} else {
		dowhat = CREATE
	}

	return dowhat
}
