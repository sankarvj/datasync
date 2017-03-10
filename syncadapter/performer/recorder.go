package performer

import (
	"database/sql"
	"log"
	"strconv"
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

func updateKey(db *sql.DB, tablename string, key int64, id int64, updated int64) {
	alreadyaddedid, localpresent := localkey(db, tablename, key)
	if localpresent { // Key Already Set
		stmt, err := db.Prepare("update " + tablename + " set synced= 'true',updated = " + strconv.FormatInt(updated, 10) + " where id=?")
		defer stmt.Close()
		if err != nil {
			log.Println("Error Prepare updating key ", err)
		}

		_, err = stmt.Exec(alreadyaddedid)
		if err != nil {
			panic(err)
		}
	} else {
		stmt, err := db.Prepare("update " + tablename + " set key=?,synced= 'true',updated = " + strconv.FormatInt(updated, 10) + " where id=?")
		defer stmt.Close()
		if err != nil {
			log.Println("Error prepare updating key ", err)
		}
		_, err = stmt.Exec(key, id)
		if err != nil {
			log.Println("Error exec updating key ", err)
		}
	}
}

func localkey(db *sql.DB, tablename string, serverid int64) (int64, bool) {
	localpresent := false
	if serverid == 0 {
		return 0, localpresent
	}

	var id int64
	sql_readall := `
	SELECT Id FROM ` + tablename + `
	WHERE Key = ` + strconv.FormatInt(serverid, 10) + ` LIMIT 1` + `
	`
	rows, err := db.Query(sql_readall)
	defer closeRows(rows)
	if err != nil {
		log.Println("Error sql readall ", sql_readall, err)
	}

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Println("Error scan ", err)
			return 0, false
		}
	}

	if id != 0 {
		localpresent = true
	}

	return id, localpresent
}

func ScanFrozenData(db *sql.DB, tablename string) bool {
	sql_readall := `
	SELECT Id FROM ` + tablename + `
	WHERE Synced = 0
	`
	rows, err := db.Query(sql_readall)
	defer closeRows(rows)
	if err != nil {
		log.Println("Error reading scanFrozenData ", err)
		return false
	}

	return rows.Next()
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}
