package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type remember struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	Date        string `json:"time"`
	ExpiredTime string `json:"expTime"`
}

func main() {
	app := fiber.New()

	remembers := []remember{}

	app.Get("/AllRemembers", AllRemembers)
	app.Post("/AddRemember", func(c *fiber.Ctx) error {
		remember := &remember{}

		if err := c.BodyParser(remember); err != nil {
			return err
		}

		if remember.Text == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Text Required"})
		}

		remember.ID = len(remembers) + 1
		remembers = append(remembers, *remember)

		return c.Status(201).JSON(remember)
	})

	log.Fatal(app.Listen(":3000"))
}

func AllRemembers(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"message": "Hello go"})
}
