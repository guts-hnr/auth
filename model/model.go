package model

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type UserRegister struct {
	Name     string `json:"name" validate:"required,min=2,max=20"`
	Age      int    `json:"age" validate:"gte=18"`
	Username string `json:"username" validate:"required,min=6,max=15"`
	Email    string `json:"email" validate:"required,min=6,max=25,email"`
	Password string `json:"password" validate:"required,min=7,max=15,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789"`
}

type UserLogin struct {
	Identifier     string `json:"identifier" validate:"required,min=5,max=25"`
	IdentifierType string `json:"identifierType" validate:"required,min=5,max=8"`
	Password       string `json:"password" validate:"required,min=7,max=15,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ProfileResponse struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func RegisterValidate(user UserRegister) error {
	return validate.Struct(user)
}

func LoginValidate(user UserLogin) error {
	return validate.Struct(user)
}
