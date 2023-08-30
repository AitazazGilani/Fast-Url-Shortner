package main

import (
	"log"
	"os"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/AitazazGilani/Fast-Url-Shortner/backend/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment file.")
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	setupRoutes(app)

	fmt.Println("Server running at port " + os.Getenv("SERVER_PORT"))
	app.Listen(os.Getenv("SERVER_PORT")) // + os.Getenv("APP_PORT"))
	// base62EncodedString := helpers.Base62Encode(9999999)
	// fmt.Println(base62EncodedString)
	// fmt.Println(helpers.Base62Decode(base62EncodedString))
}