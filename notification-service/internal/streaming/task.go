package streaming

import (
	"github.com/rfashwall/notification-service/internal/utils"
	"go.uber.org/zap"
)

type TaskEventHandler struct {
	notifier utils.Notifier
	logger   *zap.Logger
}

func NewTaskEventHandler(notifier utils.Notifier, logger *zap.Logger) *TaskEventHandler {
	return &TaskEventHandler{
		notifier: notifier,
		logger:   logger,
	}
}

func (h *TaskEventHandler) HandleEvent(event map[string]interface{}) error {
	switch event["action"] {
	case "task_created":
		h.logger.Info("Task created event received")
		return h.notifier.Send("", "userID")
	case "task_assigned":
		h.logger.Info("Task assigned event received")
		return h.notifier.Send("", "userID")
	default:
		h.logger.Warn("Unknown event action", zap.Any("event", event))
	}

	return nil
}
