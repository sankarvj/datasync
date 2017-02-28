package adapter

import (
	"database/sql"
	"log"
	"reflect"
	"strings"
)

//Expects a function followed by its params
//One of the param must implement shaper interface. This is mandatory.
//Shaper interface is responsible for datasync. Datasync skips if no params implements shaper
func CreateLocal(fn interface{}, params ...interface{}) {
	//Shaper as a rescuer
	shaperinterface := reflect.TypeOf((*Shaper)(nil)).Elem()
	var shaper Shaper

	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(params) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(params))

	for k, in := range params {
		inputs[k] = reflect.ValueOf(in)

		if reflect.TypeOf(in).Implements(shaperinterface) {
			shaper = in.(Shaper)
			tablename := strings.ToLower(reflect.TypeOf(in).Elem().Name() + "s")
			//For own table sets rt and rk as the table name
			shaper.PushReferences(tablename, tablename, tablename)

			//Form references using tags
			objtype := reflect.TypeOf(in).Elem()
			noOfFields := objtype.NumField()
			var reference_table string
			var reference_key string
			for i := 0; i < noOfFields; i++ {
				field := objtype.Field(i)
				reference_table = field.Tag.Get("rt")
				reference_key = field.Tag.Get("rk")
				if reference_table != "" && reference_key != "" {
					shaper.PushReferences(reference_table, reference_key, field.Name)
				}
			}
		}
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
	insertedId := result[0].Interface().(int64)
	if shaper != nil {
		//update the value in the local object
		shaper.UpdateLocalId(insertedId)
	}

}

func ProcessForRemote(db *sql.DB, shaper Shaper) {
	log.Println("Process for remote ", shaper.GetReferences())
	shaper.MarkAsServer(db)
}
