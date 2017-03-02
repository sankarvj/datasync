package adapter

import (
	"database/sql"
	"log"
)

func serverVal(db *sql.DB, tablename string, localid string) int64 {
	if localid == "0" {
		return 0
	}

	var key int64
	sql_readall := `
	SELECT key FROM ` + tablename + `
	WHERE id  = ` + localid + ` LIMIT 1` + `
	`

	rows, err := db.Query(sql_readall)
	defer closeRows(rows)

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

func localVal(db *sql.DB, tablename string, colname string, localid string) string {
	if localid == "0" {
		return "0"
	}

	var key string
	sql_readall := `
	SELECT ` + colname + ` FROM ` + tablename + `
	WHERE id = ` + localid + ` LIMIT 1` + `
	`
	rows, err := db.Query(sql_readall)
	defer closeRows(rows)

	if err != nil {
		log.Println("error while reading localcolid for local id", err)
		return "0"
	}

	for rows.Next() {
		err = rows.Scan(&key)
		if err != nil {
			log.Println("error while scaning localcolid for local id", err)
			key = "0"
		}
	}

	return key
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}
