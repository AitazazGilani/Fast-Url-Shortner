package routes

import(
	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
	"github.com/AitazazGilani/Fast-Url-Shortner/backend/model"
	"github.com/AitazazGilani/Fast-Url-Shortner/backend/middleware"
)


func ResolveURL(c *fiber.Ctx){



	url := c.Params("url")

	r := cache.CreateClient(0)
	defer r.Close()

	value, err := r.Get(cache.Ctx, url).Result()

	if err == redis.Nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":"short url not found in cache"
		})
	} 
	else if err != nil{
		return c.Status(fiber.StatusInternalError).JSON(fiber.Map{
			"error":"Could not connect to the Cache"
		})
	}

	rInr := cache.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(cache.Ctx, "counter")

	//redirect user
	return c.Redirect(value,301)
}