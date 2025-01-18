package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = make([]*User, 0)

func main() {
	app := fiber.New()

	// Create a new user
	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		// Generate a UUID for the user
		user.ID = uuid.New().String()

		users = append(users, user)

		return c.Status(fiber.StatusCreated).JSON(user)
	})

	// Get all users
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(users)
	})

	// Get a user by ID
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")

		for _, u := range users {
			if u.ID == userID {
				return c.JSON(u)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	})

	// Update a user by ID
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		updatedUser := new(User)

		if err := c.BodyParser(updatedUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		for i, u := range users {
			if u.ID == userID {
				users[i] = updatedUser
				users[i].ID = userID
				return c.JSON(users[i])
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	})

	// Delete a user by ID
	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")

		for i, u := range users {
			if u.ID == userID {
				users = append(users[:i], users[i+1:]...)
				return c.SendStatus(fiber.StatusNoContent)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	})

	// Start the server
	app.Listen(":3000")
}
