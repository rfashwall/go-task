package query

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rfashwall/task-service/internal/models"
)

type TaskQuery interface {
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	ListTasksByUserID(ctx context.Context, userID int) ([]*models.Task, error)
}

type MySQLTaskQuery struct {
	Conn *sql.DB
}

func NewMySQLTaskQuery(conn *sql.DB) *MySQLTaskQuery {
	return &MySQLTaskQuery{Conn: conn}
}

func (q *MySQLTaskQuery) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task
	err := q.Conn.QueryRowContext(ctx, "SELECT id, user_id, title, description, status FROM tasks WHERE id=?", id).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (q *MySQLTaskQuery) ListTasksByUserID(ctx context.Context, userID int) ([]*models.Task, error) {
	rows, err := q.Conn.QueryContext(ctx, "SELECT id, user_id, title, description, status FROM tasks WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
