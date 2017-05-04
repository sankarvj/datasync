package technique

import (
	"database/sql"
)

type Erasesync struct {
	DBInst     *sql.DB
	Tablenames []string
}

// If the user enables this technique it will check for rotten data periodically and it will automatically call the below
// methods based on the level it screwed up.

//When you completly screwed up
func DropDB() {

}

//When you partially screwed up
func DeleteTable(tablename string) {

}

//When the data stored in frozen state for a very long even after many attempt to sync that data
func EliminateRottenData(tablename string) {

}
