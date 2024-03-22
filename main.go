package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/template/handlebars/v2"

	"github.com/mshore-dev/dendrite-invite/config"
	"github.com/mshore-dev/dendrite-invite/database"
	"github.com/mshore-dev/dendrite-invite/routes"
)

func main() {
	log.Println("Hello, World!")

	engine := handlebars.New("./views/", ".hbs")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// app.Use(logger.New(logger.Config{
	// 	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	// }))

	// ratelimiting!
	// TODO: customizable ratelimit bypass list?
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Render("error", fiber.Map{"error": "Too many requests in a small time period. Chill out."})
		},
	}))

	config.LoadConfig()
	database.InitDB()

	routes.RegisterRoutes(app)

	app.Static("/static/", "./static/")

	app.Get("/ip", func(c fiber.Ctx) error {

		c.SendString("Your IP is " + c.IP())

		return nil
	})

	// app.Delete("/invite/:code", func(c fiber.Ctx) error {
	// 	return nil
	// })

	// app.Post("/invite", func(c fiber.Ctx) error {
	// 	return nil
	// })

	// app.Get("/invite/list", func(c fiber.Ctx) error {
	// 	return nil
	// })

	// app.Get("/admin", func(c fiber.Ctx) error {
	// 	return nil
	// })

	// app.Get("/admin/invites", func(c fiber.Ctx) error {
	// 	return nil
	// })

	// app.Get("/admin/logs", func(c fiber.Ctx) error {
	// 	return nil
	// })

	// TODO: configuration
	log.Fatal(app.Listen(":8080"))
}
