package streaming

import (
	"context"

	"github.com/rfashwall/task-service/internal/models"
	"github.com/rfashwall/task-service/internal/repository/command"
	"go.uber.org/zap"
)

type UserEventHandler struct {
	taskCommand command.TaskCommand
	logger      *zap.Logger
}

func NewUserEventHandler(taskCommand command.TaskCommand, logger *zap.Logger) *UserEventHandler {
	return &UserEventHandler{
		taskCommand: taskCommand,
		logger:      logger,
	}
}

func (h *UserEventHandler) HandleEvent(event map[string]interface{}) error {
	switch event["action"] {
	case "user_created":
		h.logger.Info("User created event received")
		userID := int(event["user_id"].(float64))
		h.logger.Info("Creating onboarding task", zap.Int("user_id", userID), zap.Any("event", event))
		task := &models.Task{
			UserID:      userID,
			Title:       "Onboarding Task",
			Description: "Complete your onboarding process",
			Status:      "Pending",
			AssigneeID:  userID,
		}
		h.logger.Info("Creating onboarding task", zap.Int("user_id", userID), zap.Any("task", task))
		return h.taskCommand.CreateTask(context.Background(), task)
	default:
		h.logger.Warn("Unknown event action", zap.Any("event", event))
	}

	return nil
}
