package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rfashwall/user-service/internal/repository/query"
	"go.uber.org/zap"
)

type UserQueryHandler struct {
	UserQuery query.UserQuery
	logger    *zap.Logger
}

func NewUserQueryHandler(userQuery query.UserQuery, logger *zap.Logger) *UserQueryHandler {
	return &UserQueryHandler{
		UserQuery: userQuery,
		logger:    logger,
	}
}

func (h *UserQueryHandler) SetupRoutes(app *fiber.App) {
	app.Get("/users/:id", h.getUser)
	app.Get("/users", h.listUsers)
}

func (h *UserQueryHandler) getUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user, err := h.UserQuery.GetUserByID(c.UserContext(), id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err), zap.Int("user_id", id))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(user)
}

func (h *UserQueryHandler) listUsers(c *fiber.Ctx) error {
	users, err := h.UserQuery.ListUsers(c.UserContext())
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	h.logger.Info("Listed users", zap.Int("count", len(users)))
	return c.JSON(users)
}
