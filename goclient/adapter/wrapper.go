package adapter

import (
	"log"
	"reflect"
)

//CreateLogic expects a function and its params
//One of the param must implement shaper interface
//Shaper interface is responsible for datasync. Datasync skips if no params implements shaper
func CreateLogic(fn interface{}, params ...interface{}) {
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

			//Form baseids using tags
			objtype := reflect.TypeOf(in).Elem()
			noOfFields := objtype.NumField()
			var reference_table string
			var reference_key string
			for i := 0; i < noOfFields; i++ {
				field := objtype.Field(i)
				reference_table = field.Tag.Get("rt")
				reference_key = field.Tag.Get("rk")
				shaper.PushReferences(reference_table, reference_key)
			}
		}
	}
	if shaper != nil {
		//Update synced and update values in the local object
		shaper.MarkAsLocal()
	} else {
		log.Println("datasync skips if no params implements shaper")
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
