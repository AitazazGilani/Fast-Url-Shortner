package routes

import(
	"time"
	"os"
	"fmt"
	"github.com/go-redis/redis/v8"
	cache "github.com/AitazazGilani/Fast-Url-Shortner/backend/model"
	"github.com/AitazazGilani/Fast-Url-Shortner/backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"	
	"strconv"
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
	fmt.Println("Recived request for url: " + body.URL)


	//Rate limiting step, a user can only call shortenURL only 10 times in 30 mins time
	r2 := cache.CreateClient(1)
	defer r2.Close() //close client after funciton call ends

	//get user IP in cache
	val, err := r2.Get(cache.Ctx, c.IP()).Result()
	fmt.Println("Request source IP: " + c.IP())
	if err != nil{
		fmt.Println(err.Error())
	}
	
	if err == redis.Nil{ //if not in db then set IP with quota
		_ = r2.Set(cache.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else{ //if user is in the db
		//val, _ = r2.Get(cache.Ctx, c.IP()).Result()  //WE are getting the quota the user has left in db and storing it in valInt
		valInt, _ := strconv.Atoi(val) //convert string to int

		if valInt <= 0{
			limit,_ := r2.TTL(cache.Ctx,c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":"Rate limit exceeded!",
				"rate_limit_reset": limit/time.Nanosecond/time.Second,
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
			"error": "Custom URL short is already taken",
		})
	}


	if body.Expiry == 0{
		body.Expiry = 24
	}

	err = r.Set(cache.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Unable to connect to server",
		})
	}

	defaultAPIQuotaStr := os.Getenv("API_QUOTA")
	defaultApiQuota, _ := strconv.Atoi(defaultAPIQuotaStr)
	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  defaultApiQuota,
		XRateLimitReset: 30,
	}

	//decrement db for the rate limit
	r2.Decr(cache.Ctx, c.IP())

	val, _ = r2.Get(cache.Ctx, c.IP()).Result()
	resp.XRateRemaining,_ = strconv.Atoi(val)

	ttl, _ := r2.TTL(cache.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl/time.Nanosecond/time.Minute
	
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id


	return c.Status(fiber.StatusOK).JSON(resp)
}