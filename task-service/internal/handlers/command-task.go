package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/rfashwall/task-service/internal/models"
	"github.com/rfashwall/task-service/internal/repository/command"
	"go.uber.org/zap"
)

type TaskCommandHandler struct {
	TaskCommand command.TaskCommand
	logger      *zap.Logger
	nc          *nats.Conn
}

func NewTaskCommandHandler(taskCommand command.TaskCommand, nc *nats.Conn, logger *zap.Logger) *TaskCommandHandler {
	return &TaskCommandHandler{
		TaskCommand: taskCommand,
		logger:      logger,
		nc:          nc,
	}
}

func (h *TaskCommandHandler) SetupRoutes(app *fiber.App) {
	app.Post("/tasks", h.createTask)
	app.Put("/tasks/:id", h.updateTask)
	app.Delete("/tasks/:id", h.deleteTask)

	app.Put("/tasks/:id/assign", h.assignTask)
}

func (h *TaskCommandHandler) createTask(c *fiber.Ctx) error {
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	err := h.TaskCommand.CreateTask(c.UserContext(), task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	eventData := map[string]interface{}{
		"user_id": task.UserID,
		"action":  "task_created",
	}
	eventBytes, _ := json.Marshal(eventData)
	err = h.nc.Publish("task.events", eventBytes)
	if err != nil {
		h.logger.Error("Failed to publish task created event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *TaskCommandHandler) updateTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}
	task.ID = id

	err = h.TaskCommand.UpdateTask(c.Context(), task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(task)
}

func (h *TaskCommandHandler) deleteTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	err = h.TaskCommand.DeleteTask(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaskCommandHandler) assignTask(c *fiber.Ctx) error {
	taskID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	type AssignTaskRequest struct {
		AssigneeID int `json:"assignee_id"`
	}
	var req AssignTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	// ... (Potentially add validation to check if assignee exists)

	err = h.TaskCommand.AssignTask(c.Context(), taskID, req.AssigneeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	eventData := map[string]interface{}{
		"user_id": req.AssigneeID,
		"action":  "task_assigned",
	}
	eventBytes, _ := json.Marshal(eventData)
	err = h.nc.Publish("task.events", eventBytes)
	if err != nil {
		h.logger.Error("Failed to publish task assigned event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
