package api

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"

	"github.com/mshore-dev/dendrite-invite/database"
)

// POST /api/invite/create
// create a new invite with the desired parameters.
// returns newly created invite
func apiInviteCreate(c fiber.Ctx) error {

	var inv database.Invite

	err := json.Unmarshal(c.Body(), &inv)
	if err != nil {
		c.JSON(fiber.Map{
			"success": false,
			"message": "could not unmarshal request json",
		})
		return nil
	}

	createdInv, err := database.CreateInvite(inv)
	if err != nil {
		c.JSON(fiber.Map{
			"success": false,
			"message": "could not create invite",
		})
		return nil
	}

	c.JSON(fiber.Map{
		"success": true,
		"message": "created invite",
		"invite":  createdInv,
	})

	return nil
}

// GET /api/invite/:type/:identifier
// get information about an invite, using it's code or db id
// returns invite, if found
func apiInviteGet(c fiber.Ctx) error {

	return nil
}

// PUT /api/invite/:type/:identidier
func apiInviteUpdate(c fiber.Ctx) error {

	return nil
}
