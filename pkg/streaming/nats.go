package streaming

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type EventHandler interface {
	HandleEvent(event map[string]interface{}) error
}

type NATSSubscriber struct {
	nc           *nats.Conn
	eventHandler EventHandler
	logger       *zap.Logger
}

func NewNATSSubscriber(nc *nats.Conn, eventHandler EventHandler, logger *zap.Logger) *NATSSubscriber {
	return &NATSSubscriber{
		nc:           nc,
		eventHandler: eventHandler,
		logger:       logger,
	}
}

func (s *NATSSubscriber) Subscribe(subject string) error {
	_, err := s.nc.Subscribe(subject, func(m *nats.Msg) {
		var event map[string]interface{}
		if err := json.Unmarshal(m.Data, &event); err != nil {
			s.logger.Error("Failed to unmarshal event", zap.Error(err))
			return
		}

		s.logger.Info("Received event", zap.Any("event", event))

		if err := s.eventHandler.HandleEvent(event); err != nil {
			s.logger.Error("Failed to handle event", zap.Error(err))
		}
	})
	return err
}
