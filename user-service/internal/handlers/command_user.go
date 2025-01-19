package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rfashwall/user-service/internal/models"
	"github.com/rfashwall/user-service/internal/repository/command"
	"go.uber.org/zap"
)

type UserCommandHandler struct {
	UserCommand command.UserCommand
	logger      *zap.Logger
}

func NewUserCommandHandler(userCommand command.UserCommand, logger *zap.Logger) *UserCommandHandler {
	return &UserCommandHandler{
		UserCommand: userCommand,
		logger:      logger,
	}
}

func (h *UserCommandHandler) SetupRoutes(app *fiber.App) {
	app.Post("/users", h.createUser)
	app.Put("/users/:id", h.updateUser)
	app.Delete("/users/:id", h.deleteUser)
}

func (h *UserCommandHandler) createUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		h.logger.Error("Invalid input", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	err := h.UserCommand.CreateUser(c.Context(), user)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserCommandHandler) updateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		h.logger.Error("Invalid user input", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}
	user.ID = id

	err = h.UserCommand.UpdateUser(c.Context(), user)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Error(err), zap.Int("user_id", id))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	h.logger.Info("User updated", zap.Int("id", id))
	return c.JSON(user)
}

func (h *UserCommandHandler) deleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	err = h.UserCommand.DeleteUser(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err), zap.Int("user_id", id))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	h.logger.Info("User deleted", zap.Int("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}
