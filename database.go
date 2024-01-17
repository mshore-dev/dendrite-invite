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
	CreatedAt   int64
	InviteCode  string
	ExpireTime  int64
	ExpireUses  int
	CurrentUses int
	CreatedBy   string
	Active      bool
}

type logItem struct {
	ID        int
	InviteID  string
	UserID    string
	Timestamp int64
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

func getInviteByCode(code string) (invite, error) {
	row := db.QueryRow("select * from invites where invitecode = ?", code)

	var inv invite
	err := row.Scan(&inv.ID, &inv.CreatedAt, &inv.InviteCode, &inv.ExpireTime, &inv.ExpireUses, &inv.CurrentUses, &inv.CreatedBy, &inv.Active)
	if err != nil {
		return invite{}, err
	}

	return inv, nil
}

func getInviteByID(id int) (invite, error) {
	row := db.QueryRow("select * from invites where id = ?", id)

	var inv invite
	err := row.Scan(&inv.ID, &inv.CreatedAt, &inv.InviteCode, &inv.ExpireTime, &inv.ExpireUses, &inv.CurrentUses, &inv.CreatedBy, &inv.Active)
	if err != nil {
		return invite{}, err
	}

	return inv, nil
}

// func getAllInvites() {

// }

func createInvite(inv invite) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into invites (createdat, invitecode, expiretime, expireuses, currentuses, createdby, active) values (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(time.Now().Unix(), inv.InviteCode, inv.ExpireTime, inv.ExpireUses, 0, inv.CreatedBy, 1)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func incrementInviteUses(id int) error {

	inv, err := getInviteByID(id)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update invites set currentuses = ? where id = ?")
	if err != nil {
		return err
	}

	stmt.Exec(inv.CurrentUses+1, inv.ID)
	err = tx.Commit()
	if err != nil {
		return err
	}

	// automatically check if it expires
	// useful to "use count" expires, although not time-based
	_, err = checkInviteExpires(inv.ID)
	if err != nil {
		return err
	}

	return nil
}

func checkInviteExpires(id int) (bool, error) {
	row := db.QueryRow("select * from invites where id = ?", id)

	var inv invite
	err := row.Scan(&inv.ID, &inv.CreatedAt, &inv.InviteCode, &inv.ExpireTime, &inv.ExpireUses, &inv.CurrentUses, &inv.CreatedBy, &inv.Active)
	if err != nil {
		return false, err
	}

	// ExpireTime will be set to -1 if it does not expire at a certain time
	if inv.ExpireTime != -1 && time.Now().Unix() >= inv.ExpireTime {
		// expire
		return true, setInviteActive(id, false)
	}

	if inv.CurrentUses >= inv.ExpireUses {
		// expire
		return true, setInviteActive(id, false)
	}

	return false, nil
}

func setInviteActive(id int, active bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update invites set active = ? where id = ?")
	if err != nil {
		return err
	}

	stmt.Exec(active, id)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

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
