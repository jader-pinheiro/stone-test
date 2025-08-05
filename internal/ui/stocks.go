package ui

import (
	"stone-test/internal/infra/data"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetTicker(db *gorm.DB) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ticker := ctx.Params("ticker")
		startDateStr := ctx.Query("startDate")

		if ticker == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": fiber.StatusBadRequest,
				"msg":    "É necessário informar o ticker",
			})
		}

		var startDatePtr *time.Time
		if startDateStr != "" {
			startDateParsed, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": fiber.StatusBadRequest,
					"msg":    "Formato de data inválido. Use YYYY-MM-DD (ISO 8601)",
				})
			}
			startDatePtr = &startDateParsed
		}

		data, err := data.GetTickerData(db, ctx.Context(), ticker, startDatePtr)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": fiber.StatusNotFound,
				"msg":    err.Error(),
			})
		}

		return ctx.JSON(data)
	}
}
