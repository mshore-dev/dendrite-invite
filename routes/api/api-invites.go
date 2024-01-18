package api

import (
	"encoding/json"
	"strconv"

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

	var inv database.Invite
	var err error
	var id int

	if c.Params("type") == "code" {
		inv, err = database.GetInviteByCode(c.Params("identifier"))
	} else {
		id, err = strconv.Atoi(c.Params("identifier"))
		if err != nil {
			c.JSON(fiber.Map{
				"success": false,
				"message": "could not convert id to int. is it a valid number?",
			})
			return nil
		}

		inv, err = database.GetInviteByID(id)
	}

	// fall-through(?) error handling. sounds like a bad idea lol
	if err != nil {
		c.JSON(fiber.Map{
			"success": false,
			"message": "could not look up invite",
		})
	}

	c.JSON(fiber.Map{
		"success": true,
		"message": "found invite",
		"invite":  inv,
	})
	return nil
}

// PUT /api/invite/:type/:identidier
func apiInviteUpdate(c fiber.Ctx) error {

	// TODO: code block repeated from above. could this be moved to a function?
	var inv database.Invite
	var err error
	var id int

	if c.Params("type") == "code" {
		inv, err = database.GetInviteByCode(c.Params("identifier"))
	} else {
		id, err = strconv.Atoi(c.Params("identifier"))
		if err != nil {
			c.JSON(fiber.Map{
				"success": false,
				"message": "could not convert id to int. is it a valid number?",
			})
			return nil
		}

		inv, err = database.GetInviteByID(id)
	}

	// fall-through(?) error handling. sounds like a bad idea lol
	if err != nil {
		c.JSON(fiber.Map{
			"success": false,
			"message": "could not look up invite",
		})
	}

	var parsedInv database.Invite

	err = json.Unmarshal(c.Body(), &parsedInv)
	if err != nil {
		c.JSON(fiber.Map{
			"success": false,
			"message": "could not unmarshal request json",
		})
		return nil
	}

	err = database.SetInviteActive(inv.ID, parsedInv.Active)
	if err != nil {
		c.JSON(fiber.Map{
			"success": false,
			"message": "could not update invite",
		})
		return nil
	}

	finalInv, err := database.GetInviteByID(inv.ID)
	if err != nil {
		c.JSON(fiber.Map{
			"success": false,
			"message": "could not confirm update",
		})
		return nil
	}

	c.JSON(fiber.Map{
		"success": true,
		"message": "updated invite",
		"invite":  finalInv,
	})
	return nil
}
