package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/lithammer/shortuuid/v4"
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

func GetInviteByCode(code string) (Invite, error) {
	row := db.QueryRow("select * from invites where invitecode = ?", code)

	var inv Invite
	err := row.Scan(&inv.ID, &inv.CreatedAt, &inv.InviteCode, &inv.ExpireTime, &inv.ExpireUses, &inv.CurrentUses, &inv.CreatedBy, &inv.Active)
	if err != nil {
		return Invite{}, err
	}

	return inv, nil
}

func GetInviteByID(id int) (Invite, error) {
	row := db.QueryRow("select * from invites where id = ?", id)

	var inv Invite
	err := row.Scan(&inv.ID, &inv.CreatedAt, &inv.InviteCode, &inv.ExpireTime, &inv.ExpireUses, &inv.CurrentUses, &inv.CreatedBy, &inv.Active)
	if err != nil {
		return Invite{}, err
	}

	return inv, nil
}

// func getAllInvites() {

// }

func CreateInvite(inv Invite) (Invite, error) {

	if inv.InviteCode == "" {
		log.Println("invitecode empty")
		inv.InviteCode = shortuuid.New()
	}

	inv.Active = true
	inv.CreatedAt = time.Now().Unix()

	tx, err := db.Begin()
	if err != nil {
		return Invite{}, err
	}

	stmt, err := tx.Prepare("insert into invites (createdat, invitecode, expiretime, expireuses, currentuses, createdby, active) values (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return Invite{}, err
	}

	_, err = stmt.Exec(inv.CreatedAt, inv.InviteCode, inv.ExpireTime, inv.ExpireUses, 0, inv.CreatedBy, 1)
	if err != nil {
		return Invite{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Invite{}, err
	}

	return inv, nil
}

func IncrementInviteUses(id int) error {

	inv, err := GetInviteByID(id)
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
	_, err = CheckInviteExpires(inv.ID)
	if err != nil {
		return err
	}

	return nil
}

func CheckInviteExpires(id int) (bool, error) {
	row := db.QueryRow("select * from invites where id = ?", id)

	var inv Invite
	err := row.Scan(&inv.ID, &inv.CreatedAt, &inv.InviteCode, &inv.ExpireTime, &inv.ExpireUses, &inv.CurrentUses, &inv.CreatedBy, &inv.Active)
	if err != nil {
		return false, err
	}

	// ExpireTime will be set to -1 if it does not expire at a certain time
	if inv.ExpireTime != -1 && time.Now().Unix() >= inv.ExpireTime {
		// expire
		return true, SetInviteActive(id, false)
	}

	if inv.CurrentUses >= inv.ExpireUses {
		// expire
		return true, SetInviteActive(id, false)
	}

	return false, nil
}

func SetInviteActive(id int, active bool) error {
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
