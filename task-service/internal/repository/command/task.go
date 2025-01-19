package command

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rfashwall/task-service/internal/models"
)

type TaskCommand interface {
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id int) error
	DeleteTasksByUserID(ctx context.Context, userID int) error

	AssignTask(ctx context.Context, taskID, assigneeID int) error
}

type MySQLTaskCommand struct {
	Conn *sql.DB
}

func NewMySQLTaskCommand(conn *sql.DB) *MySQLTaskCommand {
	return &MySQLTaskCommand{Conn: conn}
}

func (c *MySQLTaskCommand) CreateTask(ctx context.Context, task *models.Task) error {
	_, err := c.Conn.ExecContext(ctx, "INSERT INTO tasks (user_id, title, description, status) VALUES (?, ?, ?, ?)", task.UserID, task.Title, task.Description, task.Status)
	return err
}

func (c *MySQLTaskCommand) UpdateTask(ctx context.Context, task *models.Task) error {
	var currentStatus models.TaskStatus
	err := c.Conn.QueryRowContext(ctx, "SELECT status FROM tasks WHERE id = ?", task.ID).Scan(&currentStatus)
	if err != nil {
		return fmt.Errorf("failed to get current task status: %w", err)
	}

	if !isValidStatusTransition(currentStatus, task.Status) {
		return fmt.Errorf("invalid task status transition from '%s' to '%s'", currentStatus, task.Status)
	}

	_, err = c.Conn.ExecContext(ctx, "UPDATE tasks SET title=?, description=?, status=? WHERE id=?",
		task.Title, task.Description, task.Status, task.ID)
	return err
}

func (c *MySQLTaskCommand) DeleteTask(ctx context.Context, id int) error {
	_, err := c.Conn.ExecContext(ctx, "DELETE FROM tasks WHERE id=?", id)
	return err
}

func (c *MySQLTaskCommand) DeleteTasksByUserID(ctx context.Context, userID int) error {
	_, err := c.Conn.ExecContext(ctx, "DELETE FROM tasks WHERE user_id=?", userID)
	return err
}

func (c *MySQLTaskCommand) AssignTask(ctx context.Context, taskID, assigneeID int) error {
	result, err := c.Conn.ExecContext(ctx, "UPDATE tasks SET assignee_id=? WHERE id=?", assigneeID, taskID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return fmt.Errorf("task with ID %d not found", taskID)
	}
	return nil
}

// isValidStatusTransition checks if the transition between two task statuses is valid.
func isValidStatusTransition(currentStatus, newStatus models.TaskStatus) bool {
	switch currentStatus {
	case models.TaskStatusToDo:
		return newStatus == models.TaskStatusInProgress || newStatus == models.TaskStatusBlocked
	case models.TaskStatusInProgress:
		return newStatus == models.TaskStatusBlocked || newStatus == models.TaskStatusCompleted
	case models.TaskStatusBlocked:
		return newStatus == models.TaskStatusInProgress || newStatus == models.TaskStatusToDo
	case models.TaskStatusCompleted:
		return false // No transitions from Completed
	}
	return false
}
