package services

import (
	"context"
	"time"
)

type Comment struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	PostID      string    `json:"post_id"`
	CommentBody string    `json:"comment_body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CommentRequest struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	PostID      string    `json:"post_id"`
	CommentBody string    `json:"comment_body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Comment) GetAllComments() ([]*Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from comments`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Author,
			&comment.PostID,
			&comment.CommentBody,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (c *Comment) GetCommentsByPostID(id string) ([]*Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM comments WHERE post_id = $1`

	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Author,
			&comment.PostID,
			&comment.CommentBody,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (p *Comment) GetCommentById(id string) (*Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT * FROM comments WHERE id = $1
	`
	var comment Comment

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&comment.ID,
		&comment.Author,
		&comment.PostID,
		&comment.CommentBody,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// TODO: check if all fields is populated
func (c *Comment) UpdateComment(id string, body Comment) (*Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE comments
		SET
			author = $1,
			post_id = $2,
			comment_body = $3,
			updated_at = $4
		WHERE id=$5
	`

	_, err := db.ExecContext(
		ctx,
		query,
		body.Author,
		body.PostID,
		body.CommentBody,
		time.Now(),
		id,
	)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (c *Comment) DeleteComment(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM comment WHERE id=$1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Comment) CreateComment(comment Comment) (*Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO comments (author, post_id, comment_body, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning *
	`

	_, err := db.ExecContext(
		ctx,
		query,
		comment.Author,
		comment.PostID,
		comment.CommentBody,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
