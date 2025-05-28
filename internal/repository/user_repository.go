package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/josevitorrodriguess/gochat/internal/errors"
	"github.com/josevitorrodriguess/gochat/internal/models"
	"github.com/josevitorrodriguess/gochat/internal/utils/crypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	AuthenticateUser(ctx context.Context, email, pass string) (uuid.UUID, error)
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

func (ur *userRepository) AuthenticateUser(ctx context.Context, email, pass string) (uuid.UUID, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	var user models.User

	err := ur.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, errors.NewNotFoundError("user not found", err)
		}
		return uuid.Nil, errors.NewInternalError("database error", err)
	}
	ok := crypt.CheckPasswordHash(pass, user.Password)
	if !ok {
		return uuid.Nil, errors.NewUnauthorizedError("invalid credentials", err)
	}

	return user.ID, nil
}
