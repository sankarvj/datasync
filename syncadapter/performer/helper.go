package performer

import (
	"fmt"
	"gitlab.com/vjopensrc/datasync/syncadapter/core"
	"reflect"
	"strings"
)

func inImplementsCooker(in interface{}) bool {
	cookerin := reflect.TypeOf((*core.Cooker)(nil)).Elem()
	if reflect.TypeOf(in).Implements(cookerin) {
		return true
	} else {
		return false
	}
}

func PasserSlice(slice interface{}) []core.Passer {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	ret := make([]core.Passer, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface().(core.Passer)
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

func Tablename(in interface{}) string {
	return strings.ToLower(reflect.TypeOf(in).Elem().Name() + "s")
}

type SyncError struct {
	What string
	Stop bool
}

func (e SyncError) Error() string {
	return fmt.Sprintf("%v: %v", e.What, e.Stop)
}

func oops(errstr string, shouldstop bool) error {
	return SyncError{
		errstr,
		shouldstop,
	}
}
