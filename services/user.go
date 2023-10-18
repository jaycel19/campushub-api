package services

import (
	"context"
	"time"

	"github.com/jaycel19/campushub-api/util"
)

type User struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// USERS CREDENTIAL FOR LOGIN
type UserCreds struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) CreateUser(user User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning *
	`

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		hashedPassword,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM users`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Username,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *User) GetUserById(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM users WHERE username=$1`

	var user User

	row := db.QueryRowContext(ctx, query, username)
	err := row.Scan(
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) UserLogin(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM users WHERE username=$1`
	var user User
	row := db.QueryRowContext(ctx, query, username)
	err := row.Scan(
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
