package main

import (
	"fmt"
	"os"

	"urlShortener/shortener"
	"urlShortener/store"

	"github.com/gofiber/fiber/v2"
)

type URLMap struct {
	OriginalUrl string `json:"originalUrl"`
	ShortUrl    string `json:"shortUrl"`
}

func main() {
	app := fiber.New()
	store.InitializeStore()

	app.Post("/api/generate", func(c *fiber.Ctx) error {
		payload := new(URLMap)
		err := c.BodyParser(payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
		}

		generatedShortLink := shortener.GenerateShortLink(payload.OriginalUrl)
		saveErr := store.SaveUrlMapping(generatedShortLink, payload.OriginalUrl)
		if saveErr != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		return c.JSON(URLMap{
			ShortUrl: fmt.Sprintf("http://%v:%v/%v", os.Getenv("serverIp"), os.Getenv("serverPort"), generatedShortLink),
		})
	})

	app.Get("/:url", func(c *fiber.Ctx) error {
		shortLink := c.Params("url")
		originalUrl, err := store.RetrieveOriginalUrl(shortLink)
		if err != nil {
			fmt.Println(err)
			return c.SendString(err.Error())
		}

		return c.Redirect(originalUrl, fiber.StatusMovedPermanently)
	})

	app.Listen(fmt.Sprintf(":%v", os.Getenv("serverPort")))
}
