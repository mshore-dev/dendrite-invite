package invite

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {

	app.Get("/invite/:code", inviteCode)
	app.Post("/invite/:code/register", inviteCodeRegister)

}
