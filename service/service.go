package service

import (
	"auth/model"
	"auth/repository"
	"auth/utils"
	"errors"
)

type UserService struct {
	repository repository.UserRepository
	jwtSecret  string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repository: repo,
		jwtSecret:  jwtSecret,
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

func (s *UserService) Login(identifier string, identifierType string, password string) (*model.LoginResponse, error) {
	var user *model.User
	var err error

	switch identifierType {
	case "email":
		user, err = s.repository.GetByEmail(identifier)
	case "username":
		user, err = s.repository.GetByUsername(identifier)
	default:
		return nil, repository.ErrInvalidIdentifierType
	}

	if err != nil {
		return nil, err
	}

	if !utils.CheckPassword(user.Password, password) {
		return nil, repository.ErrInvalidPassword
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	responseToken := model.LoginResponse{
		Token: token,
	}

	return &responseToken, nil
}

func (s *UserService) GetProfile(userID int) (*model.User, error) {
	user, err := s.repository.GetById(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
