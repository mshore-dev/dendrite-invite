package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"

	"github.com/mshore-dev/dendrite-invite/config"
)

func RegisterRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:Authorization",
		Validator: validateApiKey,
	}))

	api.Post("/invite/create", apiInviteCreate)
	api.Get("/invite/:type/:identifier", apiInviteGet)
	api.Put("/invite/:type/:identifier", apiInviteUpdate)

}

func validateApiKey(c fiber.Ctx, key string) (bool, error) {
	return config.Config.AdminAPIKey == key, nil
}
