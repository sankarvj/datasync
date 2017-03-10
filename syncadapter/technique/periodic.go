package technique

import (
	"database/sql"
	"gitlab.com/vjopensrc/datasync/syncadapter/core"
	"gitlab.com/vjopensrc/datasync/syncadapter/performer"
)

type Periodic struct {
	DBInst *sql.DB
	Models []core.Cooker
}

func CreatePeriodic(db *sql.DB) Periodic {
	return Periodic{db, make([]core.Cooker, 0)}
}

func (g *Periodic) CheckPeriodic() {
	for i := 0; i < len(g.Models); i++ {
		unSyncedDataFound := performer.ScanFrozenData(g.DBInst, performer.Tablename(g.Models[i]))
		if unSyncedDataFound {
			g.Models[i].Signal(performer.TECHNIQUE_PERIODIC_SHOT)
		}
	}
}
