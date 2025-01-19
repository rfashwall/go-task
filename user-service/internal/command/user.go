package command

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rfashwall/user-service/internal/models"
)

type UserCommand interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type MySQLUserCommand struct {
	Conn *sql.DB
}

func NewMySQLUserCommand(conn *sql.DB) *MySQLUserCommand {
	return &MySQLUserCommand{Conn: conn}
}

func (c *MySQLUserCommand) CreateUser(ctx context.Context, user *models.User) error {
	_, err := c.Conn.ExecContext(ctx, "INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	return err
}

func (c *MySQLUserCommand) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := c.Conn.ExecContext(ctx, "UPDATE users SET name=?, email=?, password=? WHERE id=?", user.Name, user.Email, user.Password, user.ID)
	return err
}

func (c *MySQLUserCommand) DeleteUser(ctx context.Context, id int) error {
	_, err := c.Conn.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	return err
}
