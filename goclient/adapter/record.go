package adapter

import (
	"database/sql"
	"log"
	"strconv"
)

func serverid(db *sql.DB, tablename string, colname string, localid int64) int64 {
	if localid == 0 {
		return 0
	}

	var key int64
	sql_readall := `
	SELECT Key FROM ` + tablename + `
	WHERE ` + colname + ` = ` + strconv.FormatInt(localid, 10) + ` LIMIT 1` + `
	`
	rows, err := db.Query(sql_readall)
	defer rows.Close()

	if err != nil {
		log.Println("error while reading serverid for local id", err)
		return 0
	}

	for rows.Next() {
		err = rows.Scan(&key)
		if err != nil {
			log.Println("error while scaning serverid for local id", err)
			key = 0
		}
	}

	return key
}
