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
}

func CreateSpecificSyncer(db *sql.DB) Specificsync {
	return Specificsync{db, "", 0}
}

func CreateAdvSpecificSyncer(db *sql.DB, tablename string, localid int64) Specificsync {
	return Specificsync{db, tablename, localid}
}

//Expects a function followed by its params
//One of the param must implement shaper interface. This is mandatory.
//Shaper interface is responsible for datasync. Datasync skips if no params implements shaper
func (s *Specificsync) MakeLocal(fn interface{}, params ...interface{}) {
	var shaper Shaper
	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(params) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		if inImplementsShaper(in) {
			shaper = in.(Shaper)
		}
		inputs[k] = reflect.ValueOf(in)
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
func (s *Specificsync) CookForRemote(in interface{}) {
	//var shaper Shaper
	if inImplementsShaper(in) {
		//shaper = in.(Shaper)
		if s.Tablename == "" { //otherwise user might have set the tablename manually we don't need to set it
			s.Tablename = strings.ToLower(reflect.TypeOf(in).Elem().Name() + "s")
		}

		serverid := serverVal(s.DBInst, s.Tablename, strconv.FormatInt(s.Localid, 10))
		reflect.ValueOf(in).Elem().FieldByName("Id").SetInt(serverid)

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

func (s *Specificsync) CookFromRemote(in interface{}) {
	if inImplementsShaper(in) {
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
				serverid := reflect.ValueOf(in).Int()
				log.Println("reflect.ValueOf(in) ---- ", serverid)
				ref_col_local_val, _ := localkey(s.DBInst, reference_table, serverid)
				reflect.ValueOf(in).Elem().Field(i).SetInt(ref_col_local_val)
			}
		}
	} else {
		log.Println("No implementation of shaper found. Cannot convert it to local values")
	}
}

func (s *Specificsync) CoolItDown(key int64, updated int64) {
	updateKey(s.DBInst, s.Tablename, key, s.Localid, updated)
}

func (s *Specificsync) FindLocalItemIndex(serverid int64, dblistitems []Passer) int {
	for i := 0; i < len(dblistitems); i++ {
		if (dblistitems[i]).GetServerId() == serverid {
			return i
		}
	}
	return -1
}
