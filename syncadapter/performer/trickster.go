package performer

import (
	"gitlab.com/vjopensrc/datasync/syncadapter/core"
	"log"
	"reflect"
	"strconv"
)

//Pass by value
func (s *Pro) CookForRemote(in interface{}) {
	if inImplementsCooker(in) {
		if s.Tablename == "" { //otherwise user might have set the tablename manually we don't need to set it
			s.Tablename = Tablename(in)
		}

		if s.Localid == 0 {
			cooker := in.(core.Cooker)
			s.Localid = cooker.LocalId()
		}

		serverid := serverVal(s.DBInst, s.Tablename, strconv.FormatInt(s.Localid, 10))
		reflect.ValueOf(in).Elem().FieldByName("Id").SetInt(serverid)
		reflect.ValueOf(in).Elem().FieldByName("Key").SetInt(serverid)

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

func (s *Pro) CookFromRemote(in interface{}) {
	if inImplementsCooker(in) || inImplementsPasser(in) {
		reflect.ValueOf(in).Elem().FieldByName("Key").SetInt(reflect.ValueOf(in).Elem().FieldByName("Id").Int())
		reflect.ValueOf(in).Elem().FieldByName("Id").SetInt(0)
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
		log.Println("No implementation of cooker/passer found. Cannot convert it to local values")
	}
}

func (s *Pro) BuildUp(fn interface{}, params ...interface{}) (reflect.Value, []reflect.Value, core.Cooker) {
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

	log.Println("cooker :::", cooker)
	return f, inputs, cooker
}
