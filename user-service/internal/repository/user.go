package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rfashwall/user-service/internal/models"
)

// UserRepository defines the interface for interacting with the user repository.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

// MySqlUserRepository implements UserRepository using MySQL.
type MySqlUserRepository struct {
	DB *sql.DB
}

// NewMySQLUserRepository creates a new MySQLUserRepository.
func NewMySQLUserRepository(db *sql.DB) *MySqlUserRepository {
	return &MySqlUserRepository{DB: db}
}

// CreateUser inserts a new user into the database.
func (r *MySqlUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	result, err := r.DB.ExecContext(ctx, "INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	// Get the last inserted ID and update the User model
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(lastID)

	return nil
}

// GetUserByID retrieves a user from the database by ID.
func (r *MySqlUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRowContext(ctx, "SELECT id, name, email, password FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *MySqlUserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates an existing user in the database.
func (r *MySqlUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
		user.Name, user.Email, user.Password, user.ID)
	return err
}

// DeleteUser deletes a user from the database.
func (r *MySqlUserRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	return err
}
