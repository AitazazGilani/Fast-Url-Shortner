package routes

import(
	"time"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/AitazazGilani/Fast-Url-Shortner/backend/model"
)

type request struct{
	URL				string				`json:"url"`
	CustomShort		string				`json:"short"`
	Expiry			time.Duration		`json:"expiry"`
}

type response struct{
	URL				string				`json:"url"`
	CustomShort		string				`json:"short"`
	Expiry			time.Duration		`json:"expiry"`
	XRateRemaining	int					`json:"rate_limit"`
	XRateLimitReset	time.Duration		`json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error{

	body := new(request)

	if err := c.BodyParser(&body); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Could not parse JSON, recheck JSON fields"})
	}

	//Rate limiting step, a user can only call shortenURL only 10 times in 30 mins time
	r2 := cache.CreateClient(1)
	defer r2.Close() //close client after funciton call ends

	//get user IP in cache
	val, err := r2.Get(cacheCtx, c.IP()).Result()

	if err == redis.Nil{ //if not in db then set IP with quota
		_ = r2.Set(cache.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else{ //if user is in the db
		val, _ = r2.Get(cache.Ctx, c.IP().Result())  //WE are getting the quota the user has left in db and storing it in valInt
		valInt, _ := strconv.Atoi(val) //convert string to int

		if valInt <= 0{
			limit,_ := r2.TTL(cache.Ctx,c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":"Rate limit exceeded!",
				"rate_limit_reset": limit/time.Nanosecond/time.Second
			})
		}

	}




	//cehck if the input is an actual URL

	if !govalidator.IsURL(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid URL"})
	}

	//check for domain errors in URL

	if !helpers.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"Invalid URL"})
	}
	//enforce HTTPS/SSL

	body.URL = helpers.EnforceHTTP(body.URL)


	//Logic for creating a custom short

	var id string

	if body.CustomShort == ""{	//if custom short is empty, create a new id with 6 chars
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := cache.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(cache.Ctx, id).Result()
	if val != ""{	//if the custom short exists in db
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom URL short is already taken"
		})
	}


	if body.Expiry == 0{
		body.Expiry = 24
	}

	


	//decrement db
	r2.Decr(cache.Ctx, c.IP())
}