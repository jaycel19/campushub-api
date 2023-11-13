package services

import (
	"context"
	"time"
)

type Post struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Image       string    `json:"image"`
	PostContent string    `json:"post_content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostRequest struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	ImageData   []byte    `json:"image"`
	PostContent string    `json:"post_content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Post) GetAllPosts() ([]*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from posts`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Author,
			&post.Image,
			&post.PostContent,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (p *Post) GetPostById(id string) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT * FROM posts WHERE id = $1
	`
	var post Post

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&post.ID,
		&post.Author,
		&post.Image,
		&post.PostContent,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// TODO: check if all fields is populated
func (p *Post) UpdatePost(id string, body Post) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE posts
		SET
			author = $1,
			image = $2,
			post_content = $3,
			updated_at = $4
		WHERE id=$5
	`

	_, err := db.ExecContext(
		ctx,
		query,
		body.Author,
		body.Image,
		body.PostContent,
		time.Now(),
		id,
	)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (p *Post) DeletePost(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM posts WHERE id=$1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) CreatePost(post Post) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO posts (author, image, post_content, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning *
	`

	_, err := db.ExecContext(
		ctx,
		query,
		post.Author,
		post.Image,
		post.PostContent,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
