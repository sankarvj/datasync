package gomob

import (
	"gitlab.com/vjopensrc/datasync/goclient/model"
)

func ClearEveryThingFrom(tablename string) {
	model.ClearTable(tablename)
}
