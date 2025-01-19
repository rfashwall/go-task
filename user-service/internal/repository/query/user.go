package query

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rfashwall/user-service/internal/models"
)

type UserQuery interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)
}

type MySQLUserQuery struct {
	Conn *sql.DB
}

func NewMySQLUserQuery(conn *sql.DB) *MySQLUserQuery {
	return &MySQLUserQuery{Conn: conn}
}

func (q *MySQLUserQuery) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := q.Conn.QueryRowContext(ctx, "SELECT id, name, email, password FROM users WHERE id=?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (q *MySQLUserQuery) ListUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := q.Conn.QueryContext(ctx, "SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
