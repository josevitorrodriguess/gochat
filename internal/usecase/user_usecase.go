package usecase

import (
	"context"
	Errors "errors"

	"github.com/josevitorrodriguess/gochat/internal/errors"
	"github.com/josevitorrodriguess/gochat/internal/models/request"
	"github.com/josevitorrodriguess/gochat/internal/repository"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user request.UserRequest) error
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

	err = uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		if Errors.Is(err, errors.ErrUserAlreadyExists) {
			return errors.NewConflictError("user already exists", err)
		}
		return errors.NewInternalError("error creating user", err)
	}

	return nil
}
