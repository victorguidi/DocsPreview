package main

import (
	"docpreview/external/API"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := api.NewAPI(":5000")

	ws := app.Server.Use("/ws", app.UpgradeWss)
	ws.Get("/ws", websocket.New(app.HandleWss))

	app.Server.Route("/api", func(api fiber.Router) {
		api.Get("/", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})
	})

	app.Listen()
}
