package repository

import (
	"auth/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

var (
	ErrUserNotFoundByEmail    = errors.New("user not found by email")
	ErrUserNotFoundByUsername = errors.New("user not found by username")
	ErrEmailAlreadyExists     = errors.New("user with this email already exists")
	ErrUsernameAlreadyExists  = errors.New("user with this username already exists")
	ErrInvalidPassword        = errors.New("password is incorrect")
	ErrInvalidIdentifierType  = errors.New("identifier type is incorrect")
)

type PostrgresUserRepository struct {
	conn *pgx.Conn
}

func NewPostrgresUserRepository(conn *pgx.Conn) *PostrgresUserRepository {
	return &PostrgresUserRepository{
		conn: conn,
	}
}

func (r PostrgresUserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (name, age, username, password_hash, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.conn.QueryRow(
		context.Background(),
		query,
		user.Name,
		user.Age,
		user.Username,
		user.Password,
		user.Email,
	).Scan(&user.ID)

	return err
}

func (r PostrgresUserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User

	query := `
		SELECT id, name, age, username, password_hash, email
		FROM users
		WHERE email = $1
	`

	err := r.conn.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Username,
		&user.Password,
		&user.Email,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFoundByEmail
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r PostrgresUserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User

	query := `
		SELECT id, name, age, username, password_hash, email
		FROM users
		WHERE username = $1
	`

	err := r.conn.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Username,
		&user.Password,
		&user.Email,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFoundByUsername
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
