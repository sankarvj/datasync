package adapter

import (
	"database/sql"
	"log"
)

type Genericsync struct {
	DBInst     *sql.DB
	Tablenames []string
}

func CreateGenericSyncer(db *sql.DB) Genericsync {
	return Genericsync{db, make([]string, 0)}
}

func (g *Genericsync) SyncFrozenData() {
	for i := 0; i < len(g.Tablenames); i++ {
		unSyncedDataFound := scanFrozenData(g.DBInst, g.Tablenames[i])

		if unSyncedDataFound {
			log.Println("Sync table :: ", g.Tablenames[i])
		}
	}
}
