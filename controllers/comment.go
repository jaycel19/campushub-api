package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/services"
)

// GET/posts

func GetAllComments(w http.ResponseWriter, r *http.Request) {
	var comment services.Comment
	all, err := comment.GetAllComments()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}
	_ = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"comments": all})
}

func GetCommentById(w http.ResponseWriter, r *http.Request) {
	var commentId services.Comment
	err := json.NewDecoder(r.Body).Decode(&commentId)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	comment, err := models.Comment.GetCommentById(commentId.ID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	// TODO: handle error
	_ = helpers.WriteJSON(w, http.StatusOK, comment)
}

func GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	var postId services.Comment
	err := json.NewDecoder(r.Body).Decode(&postId)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	comments, err := models.Comment.GetCommentsByPostID(postId.PostID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	// TODO: handle error
	_ = helpers.WriteJSON(w, http.StatusOK, comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var commentReq services.CommentRequest
	err := json.NewDecoder(r.Body).Decode(&commentReq)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	commentPayload := services.Comment{
		Author:      commentReq.Author,
		PostID:      commentReq.PostID,
		CommentBody: commentReq.CommentBody,
	}
	// TODO: Handle Error
	_ = helpers.WriteJSON(w, http.StatusOK, commentPayload)
	commentCreated, err := models.Comment.CreateComment(commentPayload)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		// TODO: Handle Error
		_ = helpers.WriteJSON(w, http.StatusOK, commentCreated)
	}
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment services.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, comment)
	postObj, err := comment.UpdateComment(comment.ID, comment)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusOK, postObj)
	}
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var commentId services.Comment
	err := json.NewDecoder(r.Body).Decode(&commentId)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	err = models.Comment.DeleteComment(commentId.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_ = helpers.WriteJSON(w, http.StatusOK, "Deleted!")
}
