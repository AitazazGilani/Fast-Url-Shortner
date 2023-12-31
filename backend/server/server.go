package server

import(
	"os"
    "log"
	"github.com/AitazazGilani/Fast-Url-Shortner/backend/routes"
	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/joho/godotenv"
)

func StartServer() {
    //load env file
    err := godotenv.Load()
    
    if err != nil{
        log.Fatal("failed to load env file")
    }

    app := fiber.New()

    app.Use(logger.New())

    
    app.Get("/:url", routes.ResolveURL)
    app.Post("/api/v1", routes.ShortenURL)


    log.Fatal(app.Listen(os.Getenv("SERVER_PORT")))
}
