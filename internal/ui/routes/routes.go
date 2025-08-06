package routes

import (
	"log"
	"stone-test/internal/ui"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetRoutes(conn *gorm.DB) {

	app := fiber.New()

	app.Get("/ticker/:ticker", ui.GetTicker(conn))

	app.Get("/ticker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.StatusBadRequest,
			"msg":    "Informe um ticker na URL, ex: /ticker/PETR4",
		})
	})

	log.Fatal(app.Listen(":" + "3000"))

}
