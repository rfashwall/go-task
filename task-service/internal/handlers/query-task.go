package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rfashwall/task-service/internal/repository/query"
	"go.uber.org/zap"
)

type TaskQueryHandler struct {
	TaskQuery query.TaskQuery
	logger    *zap.Logger
}

func NewTaskQueryHandler(taskQuery query.TaskQuery, logger *zap.Logger) *TaskQueryHandler {
	return &TaskQueryHandler{
		TaskQuery: taskQuery,
		logger:    logger,
	}
}

func (h *TaskQueryHandler) SetupRoutes(app *fiber.App) {
	h.logger.Debug("setting up routes for task query service")

	app.Get("/tasks/:id", h.getTask)
	app.Get("/users/:user_id/tasks", h.listTasksByUserID)
}

func (h *TaskQueryHandler) getTask(c *fiber.Ctx) error {
	h.logger.Debug("getting task by ID")

	id, err := c.ParamsInt("id")
	if err != nil {
		h.logger.Error("invalid task ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	h.logger.Debug("fetching task by ID", zap.Int("task_id", id))
	task, err := h.TaskQuery.GetTaskByID(c.Context(), id)
	if err != nil {
		h.logger.Error("failed to fetch task by ID", zap.Error(err), zap.Int("task_id", id))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(task)
}

func (h *TaskQueryHandler) listTasksByUserID(c *fiber.Ctx) error {
	h.logger.Debug("listing tasks by user ID")

	userID, err := c.ParamsInt("user_id")
	if err != nil {
		h.logger.Error("invalid user ID", zap.Error(err), zap.String("user_id", c.Params("user_id")))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	h.logger.Debug("fetching tasks by user ID", zap.Int("user_id", userID))
	tasks, err := h.TaskQuery.ListTasksByUserID(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to fetch tasks by user ID", zap.Error(err), zap.Int("user_id", userID))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(tasks)
}
