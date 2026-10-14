package service

import (
	"auth/model"
	"auth/repository"
	"auth/utils"
	"errors"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Register(input model.UserRegister) error {
	if err := model.RegisterValidate(input); err != nil {
		return err
	}

	_, err := s.repository.GetByEmail(input.Email)

	if err == nil {
		return repository.ErrEmailAlreadyExists
	}

	if !errors.Is(err, repository.ErrUserNotFoundByEmail) {
		return err
	}

	_, err = s.repository.GetByUsername(input.Username)

	if err == nil {
		return repository.ErrUsernameAlreadyExists
	}

	if !errors.Is(err, repository.ErrUserNotFoundByUsername) {
		return err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     input.Name,
		Age:      input.Age,
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	return s.repository.Create(&user)
}

func (s *UserService) Login(identifier string, identifierType string, password string) error {
	var user *model.User
	var err error

	switch identifierType {
	case "email":
		user, err = s.repository.GetByEmail(identifier)
	case "username":
		user, err = s.repository.GetByUsername(identifier)
	default:
		return repository.ErrInvalidIdentifierType
	}

	if err != nil {
		return err
	}

	if condition := utils.CheckPassword(password, user.Password); condition != true {
		return repository.ErrInvalidPassword
	}

	return nil
}
