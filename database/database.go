package database

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

type logItem struct {
	ID        int
	InviteID  string
	UserID    string
	Timestamp int64
	IP        string
}

func InitDB() {
	// TODO: not hardcode this
	var err error
	db, err = sql.Open("sqlite", "dendrite-invite.db")
	if err != nil {
		log.Fatalf("failed to open database: %v\n", err)
	}

	createDb()
}

func createDb() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "logs" (
		"ID"	INTEGER UNIQUE,
		"InviteID"	INTEGER,
		"UserID"	TEXT,
		"Timestamp"	INTEGER,
		"IP"	TEXT,
		PRIMARY KEY("ID" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "invites" (
		"ID"	INTEGER,
		"CreatedAt"	INTEGER,
		"InviteCode"	TEXT UNIQUE,
		"ExpireTime"	INTEGER,
		"ExpireUses"	INTEGER,
		"CurrentUses"	INTEGER,
		"CreatedBy"		TEXT,
		"Active"	INTEGER,
		PRIMARY KEY("ID" AUTOINCREMENT)
	)`)

	if err != nil {
		panic(err)
	}
}
