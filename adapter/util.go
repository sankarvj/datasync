package adapter

import (
	"fmt"
	"reflect"
	"time"
)

func currentTime() int64 {
	return milliSeconds(time.Now())
}

func milliSeconds(now time.Time) int64 {
	return now.UTC().Unix() * 1000
}

func inImplementsCooker(in interface{}) bool {
	cookerin := reflect.TypeOf((*Cooker)(nil)).Elem()
	if reflect.TypeOf(in).Implements(cookerin) {
		return true
	} else {
		return false
	}
}

func PasserSlice(slice interface{}) []Passer {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	ret := make([]Passer, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface().(Passer)
	}
	return ret
}

func needUpdate(serverupdated int64, localupdated int64) bool {
	if localupdated == serverupdated {
		return false
	} else {
		return true
	}
}

type UtilError struct {
	What string
	Stop bool
}

func (e UtilError) Error() string {
	return fmt.Sprintf("%v: %v", e.What, e.Stop)
}

func oops(errstr string, shouldstop bool) error {
	return UtilError{
		errstr,
		shouldstop,
	}
}
