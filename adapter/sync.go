package adapter

import (
	"database/sql"
)

type Sync struct {
	DBInst     *sql.DB
	Tablename  string
	Localid    int64
	Chancooker chan Cooker
}
