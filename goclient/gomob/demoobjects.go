package gomob

import (
	"github.com/sankarvj/sample_syncadapter_client/goclient/model"
)

func ClearEveryThingFrom(tablename string) {
	model.ClearTable(tablename)
}
