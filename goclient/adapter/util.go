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

func inImplementsShaper(in interface{}) bool {
	//Shaper as a rescuer
	shaperinterface := reflect.TypeOf((*Shaper)(nil)).Elem()
	if reflect.TypeOf(in).Implements(shaperinterface) {
		return true
	} else {
		return false
	}

}
