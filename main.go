package main

import (
	"FiberWithGorm/database"
	"FiberWithGorm/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"digit,lowercase,uppercase,required,min=3"` // <-- a custom validation rule
}

func main() {
	database.Connect()
	app := fiber.New()
	//app.Use(logger.New())
	app.Use(cors.New())

	router.SetupRoutes(app)
	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	app.Listen(":8080")
}
