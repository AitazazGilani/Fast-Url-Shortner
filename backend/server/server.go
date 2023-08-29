package server

import (
    "fmt"
    "log"
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/joho/godotenv"
    //"github.com/streadway/amqp"
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

    shortKey := "short key logic goes here"

    response := ShortURL{
        OriginalURL: input.OriginalURL,
        ShortKey:    shortKey,
    }

    return c.JSON(response)
}

func redirectShortURL(c *fiber.Ctx) error {
    shortKey := c.Params("shortKey")
    if originalURL, ok := urlMap[shortKey]; ok {
        //logic for recording click event goes here
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
    //load env file
    err := godotenv.Load()
    
    if err != nil{
        fmt.Println(err)
    }

    app := fiber.New()

    app.Use(logger.New())

    app.Post("/shorten", shortenURL)
    app.Get("/:shortKey", redirectShortURL)

    log.Fatal(app.Listen(os.Getenv("SERVER_PORT")))
}
