package services

import (
	"database/sql"
	"time"
)

var db *sql.DB

// time for db proccess with any transaction
const dbTimeout = time.Second * 3

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{}
}

type Models struct {
	Post           Post
	PostRequest    PostRequest
	User           User
	JsonResponse   JsonResponse
	Session        Session
	Comment        Comment
	CommentRequest CommentRequest
}

type Auth struct {
	UserCreds UserCreds
}
