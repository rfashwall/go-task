package utils

import (
	"fmt"
)

type NotificationType string

const (
	Email NotificationType = "Email"
)

type Notifier interface {
	Send(recipient, message string) error
}

type EmailNotifier struct{}

func NewEmailNotifier() Notifier {
	return &EmailNotifier{}
}
func (e *EmailNotifier) Send(recipient, message string) error {
	fmt.Printf("Sending email to %s: %s\n", recipient, message)
	return nil // Simulate sending email
}
