package database

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

type Invite struct {
	ID          int    `json:"id"`
	CreatedAt   int64  `json:"created_at"`
	InviteCode  string `json:"invite_code"`
	ExpireTime  int64  `json:"expire_time"`
	ExpireUses  int    `json:"expire_uses"`
	CurrentUses int    `json:"current_uses"`
	CreatedBy   string `json:"created_by"`
	Active      bool   `json:"active"`
}

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

// func getLog() {

// }

// func getLogsByInviteID() {

// }

// func getLogsByIP() {

// }

// func getAllLogs() {

// }

// func createLog(inviteid int, ip string) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	stmt, err := tx.Prepare("insert into logs (inviteid, timestamp, ip) values (?, ?, ?)")
// 	if err != nil {
// 		return err
// 	}

// 	stmt.Exec(inviteid, time.Now().Unix(), ip)
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// 	// INSERT INTO "main"."logs" ("ID", "InviteID", "Timestamp", "IP") VALUES ('1', '', '', '');
// }
