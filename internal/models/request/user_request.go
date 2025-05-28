package request

import (
	"errors"
	"regexp"
	"strings"

	"github.com/josevitorrodriguess/gochat/internal/models"
)

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ur *UserRequest) IsValid() (bool, error) {
	if strings.TrimSpace(ur.Username) == "" {
		return false, errors.New("username is required")
	}
	if len(ur.Username) < 3 {
		return false, errors.New("username must be at least 3 characters")
	}
	if len(ur.Username) > 20 {
		return false, errors.New("username must be at most 20 characters")
	}

	if strings.TrimSpace(ur.Email) == "" {
		return false, errors.New("email is required")
	}
	if !isValidEmail(ur.Email) {
		return false, errors.New("invalid email format")
	}

	if strings.TrimSpace(ur.Password) == "" {
		return false, errors.New("password is required")
	}
	if len(ur.Password) < 6 {
		return false, errors.New("password must be at least 6 characters")
	}
	if len(ur.Password) > 100 {
		return false, errors.New("password must be at most 100 characters")
	}

	return true, nil
}

func (ur *UserRequest) ToUserModel() models.User {

	return models.User{
		Username: strings.TrimSpace(ur.Username),
		Email:    strings.TrimSpace(ur.Email),
		Password: ur.Password,
	}
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (sir *SignInRequest) IsValid() (bool, error) {
	if strings.TrimSpace(sir.Email) == "" {
		return false, errors.New("email is required")
	}
	if !isValidEmail(sir.Email) {
		return false, errors.New("invalid email format")
	}

	if strings.TrimSpace(sir.Password) == "" {
		return false, errors.New("password is required")
	}
	if len(sir.Password) < 6 {
		return false, errors.New("password must be at least 6 characters")
	}
	if len(sir.Password) > 20 {
		return false, errors.New("password must be at most 20 characters")
	}

	return true, nil
}
