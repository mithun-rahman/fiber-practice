package handler

import (
	"FiberWithGorm/model"
	"FiberWithGorm/validation"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (con *Controller) CreateUser(c *fiber.Ctx) error {
	db := con.DB
	user := new(model.User)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(&user)
	fmt.Println(user)

	if validation.IsValid(*user) == false {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": nil})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	e := db.First(&user, "username = ?", user.Username).Error
	if e == nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User exits", "data": nil})
	}

	userPassword := user.Password
	hash, erMatch := argon2id.CreateHash(userPassword, argon2id.DefaultParams)
	if erMatch != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}

	user.Password = hash
	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}
	token, err := con.GenerateToken(&model.User{
		Username: user.Username,
		Email:    user.Email,
		Address:  user.Address,
	})
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Token has created", "token": token})
}

func (con *Controller) GetAllUsers(c *fiber.Ctx) error {
	db := con.DB
	var users []model.User
	db.Find(&users)
	// If no user found, return an error
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Users not found", "data": nil})
	}
	// return users
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "Users Found", "data": users, "count": len(users)})
}

func (con *Controller) GetSingleUser(c *fiber.Ctx) error {
	db := con.DB
	// get id params
	id := c.Params("id")
	var user model.User
	// find single user in the database by id
	db.Find(&user, "id = ?", id)

	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}

func (con *Controller) UpdateUser(c *fiber.Ctx) error {
	type updateUser struct {
		Username string `json:"username"`
	}
	db := con.DB
	var user model.User
	// get id params
	id := c.Params("id")
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	user.Username = updateUserData.Username
	// Save the Changes
	db.Save(&user)
	// Return the updated user
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "users Found", "data": user})
}

func (con *Controller) DeleteUserByID(c *fiber.Ctx) error {
	db := con.DB
	var user model.User
	// get id params
	id := c.Params("id")
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	err := db.Delete(&user, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete user", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User deleted"})
}

func (con *Controller) LoginUser(c *fiber.Ctx) error {
	db := con.DB
	payload := new(model.User)
	//fromDb := new(model.User)

	err := c.BodyParser(&payload)
	username := payload.Username
	password := payload.Password
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	er := db.First(&payload, "username = ?", username).Error
	if er != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "user not exits", "data": err})
	}

	ok, e := argon2id.ComparePasswordAndHash(password, payload.Password)
	if e != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "wrong password", "data": nil})
	}
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "wrong password", "data": nil})
	}
	if ok {
		token, err := con.GenerateToken(&model.User{
			Username: payload.Username,
			Email:    payload.Email,
			Address:  payload.Address,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "something wrong", "data": nil})
		}
		return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Token has created", "token": token})
	}
	return nil
}
