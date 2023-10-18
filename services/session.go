package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    string    `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *Session) CreateSession(session Session) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO sessions (id, username, refresh_token, expires_at, created_at)
		values ($1, $2, $3, $4, $5) returning *
	`

	err := db.QueryRowContext(
		ctx,
		query,
		session.ID,
		session.Username,
		session.RefreshToken,
		time.Now().Add(time.Hour*24*30),
		time.Now(),
	).Scan(
		&session.ID,
		&session.Username,
		&session.RefreshToken,
		&session.IsBlocked,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (S *Session) GetSessionById(id string) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM sessions WHERE id=$1`

	var session Session
	fmt.Println(id)
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&session.ID,
		&session.Username,
		&session.RefreshToken,
		&session.IsBlocked,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
