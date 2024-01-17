package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/handlebars/v2"

	"log"
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

	loadConfig()
	initDb()

	app.Static("/static/", "./static/")

	app.Get("/invite/:code", func(c fiber.Ctx) error {

		inv, err := getInviteByCode(c.Params("code"))
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

		expired, err := checkInviteExpires(inv.ID)
		if err != nil {
			c.Render("error", fiber.Map{"error": "could not check if invite is expired"})
		}

		if expired {
			log.Printf("invite %s has (just) expired\n", inv.InviteCode)
			c.Render("invite-expired", fiber.Map{})
			return nil
		}

		c.Render("invite", fiber.Map{
			"instanceName": cfg.InstanceName,
			"inviteCode":   inv.InviteCode,
			"error":        "",
		})
		return nil
	})

	app.Post("/invite/:code/register", func(c fiber.Ctx) error {

		log.Printf("got request for /invite/%s/register/\n", c.Params("code"))

		inv, err := getInviteByCode(c.Params("code"))
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
				"instanceName": cfg.InstanceName,
				"inviteCode":   inv.InviteCode,
				"error":        "username must be between 4-18 characters",
			})
			return nil
		}

		if password != password2 {
			c.Render("invite", fiber.Map{
				"instanceName": cfg.InstanceName,
				"inviteCode":   inv.InviteCode,
				"error":        "passwords do not match",
			})
			return nil
		}

		if len(password) > 32 || len(password) < 8 {
			c.Render("invite", fiber.Map{
				"instanceName": cfg.InstanceName,
				"inviteCode":   inv.InviteCode,
				"error":        "password must be between 8 and 32 characters",
			})
			return nil
		}

		userid, token, err := matrixRegister(username, password)
		if err != nil {
			c.Render("error", fiber.Map{"error": "failed to register. don't try again."})
			log.Printf("failed to create account: %v\n", err)
			return nil
		}

		c.Render("success", fiber.Map{
			"instanceName": cfg.InstanceName,
			"clientUrl":    cfg.ClientURL,
			"accountId":    userid,
			"accessToken":  token,
		})

		createLog(inv.ID, c.IP())

		return nil
	})

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
