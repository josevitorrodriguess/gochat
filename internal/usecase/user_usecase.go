package usecase

import (
	"context"
	Errors "errors"

	"github.com/google/uuid"
	"github.com/josevitorrodriguess/gochat/internal/errors"
	"github.com/josevitorrodriguess/gochat/internal/models/request"
	"github.com/josevitorrodriguess/gochat/internal/repository"
	"github.com/josevitorrodriguess/gochat/internal/utils/crypt"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user request.UserRequest) error
	AuthenticateUser(ctx context.Context, email, pass string) (uuid.UUID, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (uc *userUseCase) CreateUser(ctx context.Context, userReq request.UserRequest) error {
	ok, err := userReq.IsValid()
	if err != nil {
		return errors.NewInternalError("error validating user", err)
	}
	if !ok {
		return errors.NewBadRequestError("invalid user data", nil)
	}

	user := userReq.ToUserModel()

	user.Password, err = crypt.HashPassword(user.Password)
	if err != nil {
		return errors.NewInternalError("error hashing password", err)
	}

	err = uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		if Errors.Is(err, errors.ErrUserAlreadyExists) {
			return errors.NewConflictError("user already exists", err)
		}
		return errors.NewInternalError("error creating user", err)
	}

	return nil
}

func (uc *userUseCase) AuthenticateUser(ctx context.Context, email, pass string) (uuid.UUID, error) {
	if email == "" || pass == "" {
		return uuid.Nil, errors.NewBadRequestError("email and password are required", nil)
	}

	userID, err := uc.userRepo.AuthenticateUser(ctx, email, pass)
	if err != nil {
		if Errors.Is(err, Errors.New("not found")) {
			return uuid.Nil, errors.NewNotFoundError("user not found", err)
		}
		return uuid.Nil, errors.NewInternalError("error authenticating user", err)
	}

	return userID, nil
}
