package database

import "time"

type Log struct {
	ID        int    `json:"id"`
	InviteID  int    `json:"invite_id"`
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	IP        string `json:"ip"`
}

func CreateLog(inviteId int, userId, ip string) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into logs (inviteid, userid, timestamp, ip) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(inviteId, userId, time.Now().Unix(), ip)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
