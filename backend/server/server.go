package server

import (
    "encoding/json"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/streadway/amqp"
    "net/http"
)

type ShortURL struct {
    OriginalURL string `json:"original_url"`
    ShortKey    string `json:"short_key"`
}

var urlMap = make(map[string]string)

func shortenURL(c *fiber.Ctx) error {
    var input ShortURL
    if err := c.BodyParser(&input); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON data"})
    }

    shortKey := generateShortKey()

    urlMap[shortKey] = input.OriginalURL
    storeMappingInDatabase(shortKey, input.OriginalURL)

    publishURLCreationEvent(shortKey)

    response := ShortURL{
        OriginalURL: input.OriginalURL,
        ShortKey:    shortKey,
    }

    return c.JSON(response)
}

func redirectShortURL(c *fiber.Ctx) error {
    shortKey := c.Params("shortKey")
    if originalURL, ok := urlMap[shortKey]; ok {
        recordClickEvent(shortKey)
        return c.Redirect(originalURL, http.StatusFound)
    } else {
        return c.Status(http.StatusNotFound).SendString("Short URL not found")
    }
}

func publishURLCreationEvent(shortKey string) {
    // Similar to the previous example
}

func startURLCreationEventConsumer() {
    // Similar to the previous example
}

func startServer() {
    app := fiber.New()

    app.Post("/shorten", shortenURL)
    app.Get("/:shortKey", redirectShortURL)

    go startURLCreationEventConsumer() // Start the consumer Goroutine

    err := app.Listen(":8080")
    if err != nil {
        fmt.Println("Error starting the server:", err)
    }
}
