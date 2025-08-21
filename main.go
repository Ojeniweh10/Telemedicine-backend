package main

import (
	"context"
	"log"
	"telemed/config"
	"telemed/database"
	"telemed/routes"
	"telemed/servers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	servers.Ctx = context.Background()
	servers.Db = database.NewConnection()
	app := fiber.New(fiber.Config{
		AppName: "Telemedicine Backend",
	})
	app.Use(logger.New())
	app.Get("/admin/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Welcome to Telemed Backend",
		})
	})
	app.Use(func(c *fiber.Ctx) error {
		auth := c.Get("G-auth")
		if auth == "" || auth != config.GatewaySecret {
			log.Println("Invalid password for authentication")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "access denied, invalid route authentication",
			})
		}
		return c.Next()
	})
	routes.AdminRoutes(app)
	routes.Routes(app)
	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Route not found",
		})
	})

	log.Fatal(app.Listen(":8080"))
}
