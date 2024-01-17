package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

type invite struct {
	ID          int
	CreatedAt   time.Time
	InviteCode  string
	ExpireTime  time.Time
	ExpireUses  int
	CurrentUses int
	Active      bool
}

type logItem struct {
	ID        int
	InviteID  string
	Timestamp time.Time
	IP        string
}

func initDb() {
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
		"Active"	INTEGER,
		PRIMARY KEY("ID" AUTOINCREMENT)
	)`)

	if err != nil {
		panic(err)
	}
}

func getInviteByCode(code string) (invite, error) {
	row := db.QueryRow("select id, invitecode, active from invites where invitecode = ?", code)

	var inv invite
	err := row.Scan(&inv.ID, &inv.InviteCode, &inv.Active)
	if err != nil {
		return invite{}, err
	}

	return inv, nil
}

// func getInviteByID() {

// }

// func getAllInvites() {

// }

// func createInvite() {

// }

// func expireInvite() {

// }

// func getLog() {

// }

// func getLogsByInviteID() {

// }

// func getLogsByIP() {

// }

// func getAllLogs() {

// }

func createLog(inviteid int, ip string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into logs (inviteid, timestamp, ip) values (?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(inviteid, time.Now().Unix(), ip)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
	// INSERT INTO "main"."logs" ("ID", "InviteID", "Timestamp", "IP") VALUES ('1', '', '', '');
}
