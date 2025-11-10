package usecases

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"github.com/ninahf618/go-todo-api/pkg/auth"
	"github.com/ninahf618/go-todo-api/pkg/security"
)

type UserUsecase struct {
	userRepo    repositories.UserRepository
	tokenSvc    auth.TokenService
	passwordSvc security.PasswordService
}

func NewUserUsecase(r repositories.UserRepository, t auth.TokenService, p security.PasswordService) *UserUsecase {
	return &UserUsecase{
		userRepo:    r,
		tokenSvc:    t,
		passwordSvc: p,
	}
}

func (u *UserUsecase) Register(username, email, password string) (*entities.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)
	if username == "" || email == "" || password == "" {
		return nil, errors.New("missing required fields")
	}

	existing, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := u.passwordSvc.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		ID:       uuid.NewString(),
		Username: username,
		Email:    email,
		Password: hash,
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", errors.New("missing credentials")
	}

	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil || !u.passwordSvc.Verify(user.Password, password) {
		return "", errors.New("invalid email or password")
	}

	return u.tokenSvc.Generate(user.ID)
}
