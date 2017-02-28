package adapter

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func coreReference(db *sql.DB, basemodel *BaseModel) error {
	var err error
	var mainreference Reference
	references := basemodel.GetReferences()
	for i := 0; i < len(*references); i++ {
		if isMainTableRef((*references)[i]) {
			mainreference = (*references)[i]
			(*references)[i].localkey = strconv.FormatInt(basemodel.Id, 10)
			basemodel.Id = serverid(db, (*references)[i].tablename, "id", basemodel.Id)
		}
	}

	if mainreference == (Reference{}) || mainreference.tablename == "" {
		err = oops("Table name not set", true)
		fmt.Println("Table name not set during CreateLocal call")
	} else { // This will pass only if main table available
		for i := 0; i < len(*references); i++ {
			if !isMainTableRef((*references)[i]) { //Loop all references except the main table reference

			}
		}
	}

	return err
}

func isMainTableRef(ref Reference) bool {
	return ref.colname == ref.tablename
}

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
