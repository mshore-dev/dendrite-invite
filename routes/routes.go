package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/mshore-dev/dendrite-invite/routes/api"
	"github.com/mshore-dev/dendrite-invite/routes/invite"
)

func RegisterRoutes(app *fiber.App) {

	api.RegisterRoutes(app)
	invite.RegisterRoutes(app)

}
