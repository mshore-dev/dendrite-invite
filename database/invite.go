package database

import (
	"log"
	"time"

	"github.com/lithammer/shortuuid/v4"
)

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
