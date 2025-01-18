package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rfashwall/user-service/internal/models"
	"github.com/rfashwall/user-service/internal/repository"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	Repo   repository.UserRepository
	Logger *zap.Logger
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(repo repository.UserRepository, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		Repo:   repo,
		Logger: logger,
	}
}

// CreateUser handles the POST request to create a new user.
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := &models.User{}

	if err := c.BodyParser(user); err != nil {
		h.Logger.Error("Error parsing request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	if err := h.Repo.CreateUser(c.Context(), user); err != nil {
		h.Logger.Error("Error creating user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByID handles the GET request to retrieve a user by ID.
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	user, err := h.Repo.GetUserByID(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to get user"})
	}

	return c.JSON(user)
}

// GetAllUsers handles the GET request to retrieve all users.
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.Repo.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to get users"})
	}
	return c.JSON(users)
}

// UpdateUser handles the PUT request to update an existing user.
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	user.ID = id

	if err := h.Repo.UpdateUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update user"})
	}

	return c.JSON(user)
}

// DeleteUser handles the DELETE request to delete a user.
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	if err := h.Repo.DeleteUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete user"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
