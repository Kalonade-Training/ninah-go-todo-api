package usecases

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ninahf618/go-todo-api/domain/entities"
	"github.com/ninahf618/go-todo-api/domain/repositories"
	"github.com/ninahf618/go-todo-api/infrastructure/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserUsecase struct {
	repo repositories.UserRepository
}

func NewUserUsecase(repo repositories.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Register(username, email, password string) (*entities.User, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	existingUser, err := u.repo.FindByEmail(normalizedEmail)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		ID:       uuid.NewString(),
		Username: username,
		Email:    normalizedEmail,
		Password: string(hashed),
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) Login(email, password string) (string, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	user, err := u.repo.FindByEmail(normalizedEmail)

	if err != nil || user == nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
