package service

import (
	"context"
	"errors"
	"strings"

	"github.com/deblinux123/Redis-DB/day-07/internal/model"
	"github.com/deblinux123/Redis-DB/day-07/internal/repository"
)

var (
	ErrInvlaidUser  = errors.New("Invlaid User")
	ErrUserNotFound = errors.New("User Not Found.")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, user model.User) error {
	user.Name = strings.TrimSpace(user.Name)

	if user.ID == "" || user.Name == "" || user.Age <= 0 {
		return ErrInvlaidUser
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (model.User, error) {
	if strings.TrimSpace(id) == "" {
		return model.User{}, ErrInvlaidUser
	}

	user, err := s.repo.GetUserByID(ctx, id)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
