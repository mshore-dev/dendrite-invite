package invite

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/mshore-dev/dendrite-invite/config"
	"github.com/mshore-dev/dendrite-invite/database"
	"github.com/mshore-dev/dendrite-invite/matrix"
)

func inviteCode(c fiber.Ctx) error {
	inv, err := database.GetInviteByCode(c.Params("code"))
	if err != nil {
		log.Printf("failed to find code %s: %v\n", c.Params("code"), err)
		c.Render("invite-invalid", fiber.Map{})
		return nil
	}

	log.Printf("invite: %v\n", inv)

	if !inv.Active {
		log.Printf("invite %s is marked inactive\n", inv.InviteCode)
		c.Render("invite-expired", fiber.Map{})
		return nil
	}

	expired, err := database.CheckInviteExpires(inv.ID)
	if err != nil {
		c.Render("error", fiber.Map{"error": "could not check if invite is expired"})
	}

	if expired {
		log.Printf("invite %s has (just) expired\n", inv.InviteCode)
		c.Render("invite-expired", fiber.Map{})
		return nil
	}

	c.Render("invite", fiber.Map{
		"instanceName": config.Config.InstanceName,
		"inviteCode":   inv.InviteCode,
		"inviter":      inv.CreatedBy,
		"error":        "",
	})
	return nil
}

func inviteCodeRegister(c fiber.Ctx) error {
	log.Printf("got request for /invite/%s/register/\n", c.Params("code"))

	inv, err := database.GetInviteByCode(c.Params("code"))
	if err != nil {
		log.Printf("failed to find code %s: %v\n", c.Params("code"), err)
		c.Render("error", fiber.Map{"error": "invite is invalid or has expired"})
		return nil
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	password2 := c.FormValue("password2")

	log.Println(len(username))

	if len(username) > 18 || len(username) < 4 {
		c.Render("invite", fiber.Map{
			"instanceName": config.Config.InstanceName,
			"inviteCode":   inv.InviteCode,
			"inviter":      inv.CreatedBy,
			"error":        "username must be between 4-18 characters",
		})
		return nil
	}

	if password != password2 {
		c.Render("invite", fiber.Map{
			"instanceName": config.Config.InstanceName,
			"inviteCode":   inv.InviteCode,
			"inviter":      inv.CreatedBy,
			"error":        "passwords do not match",
		})
		return nil
	}

	if len(password) > 32 || len(password) < 8 {
		c.Render("invite", fiber.Map{
			"instanceName": config.Config.InstanceName,
			"inviteCode":   inv.InviteCode,
			"inviter":      inv.CreatedBy,
			"error":        "password must be between 8 and 32 characters",
		})
		return nil
	}

	userid, token, err := matrix.Register(username, password)
	if err != nil {
		c.Render("error", fiber.Map{"error": "failed to register. don't try again."})
		log.Printf("failed to create account: %v\n", err)
		return nil
	}

	c.Render("success", fiber.Map{
		"instanceName": config.Config.InstanceName,
		"clientUrl":    config.Config.ClientURL,
		"accountId":    userid,
		"accessToken":  token,
	})

	database.IncrementInviteUses(inv.ID)

	// createLog(inv.ID, c.IP())

	return nil
}
