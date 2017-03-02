package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//**********Create Quries**********

//Ticket table create
var sql_create_ticket_table = `
	CREATE TABLE IF NOT EXISTS tickets(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Key INTEGER DEFAULT 0,
		Subject TEXT,
		Desc TEXT,
		Requester TEXT,
		Agent TEXT,
		Updated LONG,
		Created DATETIME,
		Synced BOOLEAN DEFAULT FALSE
	);
	`

//Note table create
//Relationship - Each ticket has many notes
var sql_create_note_table = `
	CREATE TABLE IF NOT EXISTS notes(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Key INTEGER DEFAULT 0,
		Ticketid INTEGER,
		Name TEXT,
		Desc TEXT,
		Updated LONG,
		Created DATETIME,
		Synced BOOLEAN DEFAULT FALSE,
		FOREIGN KEY(Ticketid) REFERENCES tickets (Id) ON DELETE CASCADE
	);
	`

//Notemeta table create
//Relationship - Each note has many notemetas
var sql_create_attachment_table = `
	CREATE TABLE IF NOT EXISTS attachments(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Key INTEGER DEFAULT 0,
		Noteid INTEGER,
		Url TEXT,
		Size INTEGER,
		Updated LONG,
		Created DATETIME,
		Synced BOOLEAN DEFAULT FALSE,
		FOREIGN KEY(Noteid) REFERENCES notes (Id) ON DELETE CASCADE
	);
	`

//**********Insert Quries**********

//ticket insert query
var sql_ticket_insert_query = `
	INSERT INTO tickets(
		Key,
		Subject,
		Desc,
		Requester,
		Agent,
		Updated,
		Created,
		Synced
	) values(?, ?, ?, ?, ?, ?, ?, ?)
	`

//note insert query
var sql_note_insert_query = `
	INSERT INTO notes(
		Key,
		Ticketid,
		Name,
		Desc,
		Updated,
		Created,
		Synced
	) values(?, ?, ?, ?, ?, ?, ?)
	`

//attchment insert query
var sql_attachment_insert_query = `
	INSERT INTO attachments(
		Key,
		Noteid,
		Url,
		Size,
		Updated,
		Created,
		Synced
	) values(?, ?, ?, ?, ?, ?, ?)
	`

var frontendAdapter FrontendAdapter

//Adapter talks with frontend and get back info needed for goclient
type FrontendAdapter struct {
}

func RegisterFrontendAdapter(f FrontendAdapter) {
	frontendAdapter = f
}

func (f FrontendAdapter) DatabasePath() string {
	return "datasync.db"
}

var db *sql.DB

func InitDB() *sql.DB {
	var err error

	if db == nil {
		db, err = sql.Open("sqlite3", frontendAdapter.DatabasePath()+"?mode=rwc")
		db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			log.Println("database err ", err)
			return nil
		}
		createTable(db)
	} else {
		log.Println("....,,.... Using old db connection")
	}
	return db
}

func createTable(db *sql.DB) {
	// create tickets table if not exists
	_, err := db.Exec(sql_create_ticket_table)
	if err != nil {
		log.Println("database table create err ", err)
		return
	}
	// create notes table if not exists
	_, err = db.Exec(sql_create_note_table)
	if err != nil {
		log.Println("database table create err ", err)
		return
	}
	// create attachments table if not exists
	_, err = db.Exec(sql_create_attachment_table)
	if err != nil {
		log.Println("database table create err ", err)
		return
	}
}

func StoreTicket(ticket *Ticket) int64 {
	db := InitDB()
	stmt, err := db.Prepare(sql_ticket_insert_query)
	defer stmt.Close()
	if err != nil {
		log.Println("database prepare insert ticket sql err ", err)
		return 0
	}

	var result sql.Result
	result, err = stmt.Exec(ticket.Id, ticket.Subject, ticket.Desc, ticket.requester, ticket.agent, ticket.Updated, ticket.created, ticket.Synced)
	if err != nil {
		log.Println("database insert ticket sql err ", err)
		return 0
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		return 0
	} else {
		return insertedId
	}
}

func StoreNote(note *Note) int64 {
	db := InitDB()
	stmt, err := db.Prepare(sql_note_insert_query)
	defer stmt.Close()
	if err != nil {
		log.Println("database prepare insert note sql err ", err)
		return 0
	}

	var result sql.Result
	result, err = stmt.Exec(note.Id, note.Ticketid, note.Name, note.Desc, note.Updated, note.created, note.Synced)
	if err != nil {
		log.Println("database insert note sql err ", err)
		return 0
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		return 0
	} else {
		return insertedId
	}
}
