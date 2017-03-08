package adapter

import (
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	NOTHING = 0
	CREATE  = 1
	UPDATE  = 2
)

type Specificsync struct {
	DBInst    *sql.DB
	Tablename string
	Localid   int64
}

func CreateSpecificSyncer(db *sql.DB) Specificsync {
	return Specificsync{db, "", 0}
}

func CreateAdvSpecificSyncer(db *sql.DB, tablename string, localid int64) Specificsync {
	return Specificsync{db, tablename, localid}
}

//Expects a function followed by its params
//One of the param must implement cooker interface. This is mandatory.
//Cooker interface is responsible for datasync. Datasync skips if no params implements cooker
func (s *Specificsync) MakeLocal(fn interface{}, params ...interface{}) {
	var cooker Cooker
	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(params) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		if inImplementsCooker(in) {
			cooker = in.(Cooker)
		}
		inputs[k] = reflect.ValueOf(in)
	}

	if cooker != nil {
		//Update synced and update values in the local object
		cooker.MarkAsLocal()
	} else {
		log.Println("datasync skips since there is no params passed to this func implements cooker!")
	}
	//call the corresponding method
	result := f.Call(inputs)
	//find the last inserted id
	s.Localid = result[0].Interface().(int64)
	if cooker != nil {
		//update the value in the local object
		cooker.UpdateLocalId(s.Localid)
	}

}

func (s *Specificsync) UpdateLocal(fn interface{}, params ...interface{}) {
	var cooker Cooker
	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(params) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		if inImplementsCooker(in) {
			cooker = in.(Cooker)
		}
		inputs[k] = reflect.ValueOf(in)
	}

	if cooker != nil {
		//Update synced and update values in the local object
		cooker.MarkAsPureLocal()
	} else {
		log.Println("datasync skips since there is no params passed to this func implements cooker!")
	}
	//call the corresponding method
	f.Call(inputs)
}

//Pass by value
func (s *Specificsync) CookForRemote(in interface{}) {
	if inImplementsCooker(in) {
		if s.Tablename == "" { //otherwise user might have set the tablename manually we don't need to set it
			s.Tablename = strings.ToLower(reflect.TypeOf(in).Elem().Name() + "s")
		}

		if s.Localid == 0 {
			passer := in.(Passer)
			s.Localid = passer.GetLSId()
		}

		serverid := serverVal(s.DBInst, s.Tablename, strconv.FormatInt(s.Localid, 10))
		reflect.ValueOf(in).Elem().FieldByName("Id").SetInt(serverid)

		//Form references using tags
		objtype := reflect.TypeOf(in).Elem()
		noOfFields := objtype.NumField()
		var reference_table string
		for i := 0; i < noOfFields; i++ {
			field := objtype.Field(i)
			reference_table = field.Tag.Get("rt")
			if reference_table != "" {
				ref_col_local_val := localVal(s.DBInst, s.Tablename, field.Name, strconv.FormatInt(s.Localid, 10))
				sercolval := serverVal(s.DBInst, reference_table, ref_col_local_val)
				reflect.ValueOf(in).Elem().Field(i).SetInt(sercolval)
			}
		}
	} else {
		log.Println("No implementation of cooker found. Cannot annex remote values")
	}
}

func (s *Specificsync) CookFromRemote(in interface{}) {
	if inImplementsCooker(in) {
		//Form references using tags
		objtype := reflect.TypeOf(in).Elem()
		noOfFields := objtype.NumField()
		var reference_table string
		var reference_key string
		for i := 0; i < noOfFields; i++ {
			field := objtype.Field(i)
			reference_table = field.Tag.Get("rt")
			reference_key = field.Tag.Get("rk") //Used here
			if reference_table != "" && reference_key != "" {
				serverid := reflect.ValueOf(in).Elem().Field(i).Int()
				ref_col_local_val, _ := localkey(s.DBInst, reference_table, serverid)
				reflect.ValueOf(in).Elem().Field(i).SetInt(ref_col_local_val)
			}
		}
	} else {
		log.Println("No implementation of cooker found. Cannot convert it to local values")
	}
}

func (s *Specificsync) CoolItDown(key int64, updated int64) {
	updateKey(s.DBInst, s.Tablename, key, s.Localid, updated)
}

func (s *Specificsync) HotId(tablename string, id int64) int64 {
	return serverVal(s.DBInst, tablename, strconv.FormatInt(id, 10))
}

func (s *Specificsync) WhatToDo(passer Passer, dblistitems []Passer) (dbprimaryid int64, dowhat int64) {
	//HOT to COLD conversion
	s.CookFromRemote(passer)
	//Though passer is cold its id param is still hot.
	serverid := passer.GetLSId()
	updated := passer.GetUpdatedAt()

	index := -1
	dowhat = NOTHING
	dbprimaryid = 0
	for i := 0; i < len(dblistitems); i++ {
		if (dblistitems[i]).GetServerId() == serverid {
			index = i
		}
	}

	if index != -1 { //already stored
		if needUpdate(updated, dblistitems[index].GetUpdatedAt()) {
			dowhat = UPDATE
			dbprimaryid = dblistitems[index].GetLSId()
		}
	} else {
		dowhat = CREATE
	}

	return dbprimaryid, dowhat
}
