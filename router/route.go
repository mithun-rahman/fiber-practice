package router

import (
	"FiberWithGorm/database"
	"FiberWithGorm/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/")
	v1 := api.Group("/user")
	// routes

	con := handler.Controller{}
	con.DB = database.DB.Db

	signBytes, _ := ioutil.ReadFile("./app.rsa")
	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	con.PrivateKey = signKey

	verifyBytes, _ := ioutil.ReadFile("./app.rsa.pub")
	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	con.PublicKey = verifyKey

	v1.Get("/", con.GetAllUsers)
	v1.Get("/:id", con.GetSingleUser)
	v1.Post("/", con.CreateUser)
	v1.Patch("/:id", con.UpdateUser)
	v1.Delete("/:id", con.DeleteUserByID)
	v1.Post("/log", con.LoginUser)
}
