package adapter

import (
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Specificsync struct {
	DBInst    *sql.DB
	Tablename string
	Localid   int64
	Serverid  int64
}

//Expects a function followed by its params
//One of the param must implement shaper interface. This is mandatory.
//Shaper interface is responsible for datasync. Datasync skips if no params implements shaper
func (s *Specificsync) CreateLocal(fn interface{}, params ...interface{}) {
	var shaper Shaper
	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(params) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		inputs[k] = reflect.ValueOf(in)
	}

	if inImplementsShaper(in) {
		shaper = in.(Shaper)
	}

	if shaper != nil {
		//Update synced and update values in the local object
		shaper.MarkAsLocal()
	} else {
		log.Println("datasync skips since there is no params passed to this func implements shaper!")
	}
	//call the corresponding method
	result := f.Call(inputs)
	//find the last inserted id
	s.Localid = result[0].Interface().(int64)
	if shaper != nil {
		//update the value in the local object
		shaper.UpdateLocalId(s.Localid)
	}

}

//Pass by value
func (s *Specificsync) AnnexRemote(in interface{}) {
	//var shaper Shaper
	if inImplementsShaper(in) {
		//shaper = in.(Shaper)
		s.Tablename = strings.ToLower(reflect.TypeOf(in).Elem().Name() + "s")
		s.Serverid = serverVal(s.DBInst, s.Tablename, strconv.FormatInt(s.Localid, 10))

		//Form references using tags
		objtype := reflect.TypeOf(in).Elem()
		noOfFields := objtype.NumField()
		var reference_table string
		var reference_key string
		for i := 0; i < noOfFields; i++ {
			field := objtype.Field(i)
			reference_table = field.Tag.Get("rt")
			reference_key = field.Tag.Get("rk") //Not used
			if reference_table != "" && reference_key != "" {
				ref_col_local_val := localVal(s.DBInst, s.Tablename, field.Name, strconv.FormatInt(s.Localid, 10))
				sercolval := serverVal(s.DBInst, reference_table, ref_col_local_val)
				reflect.ValueOf(in).Elem().Field(i).SetInt(sercolval)
			}
		}
	} else {
		log.Println("No implementation of shaper found. Cannot annex remote values")
	}
}
