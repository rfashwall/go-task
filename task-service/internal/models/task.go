package models

type Task struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`

	AssigneeID int `json:"assignee_id"`
}

// TaskStatus represents the possible statuses for a task.
type TaskStatus string

// Valid task statuses.
const (
	TaskStatusToDo       TaskStatus = "To Do"
	TaskStatusInProgress TaskStatus = "In Progress"
	TaskStatusBlocked    TaskStatus = "Blocked"
	TaskStatusCompleted  TaskStatus = "Completed"
)
