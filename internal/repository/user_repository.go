package repository

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/josevitorrodriguess/gochat/internal/errors"
	"github.com/josevitorrodriguess/gochat/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (username, email, password) VALUES (:username, :email, :password)`

	_, err := ur.db.NamedExecContext(ctx, query, user)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return errors.NewConflictError("user already exists", err)
		}
		return errors.NewInternalError("database error", err)
	}

	return nil
}
