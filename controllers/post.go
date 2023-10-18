package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/services"
)

var models services.Models

// GET/posts

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	all, err := post.GetAllPosts()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}
	_ = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"posts": all})
}

func GetPostById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := models.Post.GetPostById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	// TODO: handle error
	_ = helpers.WriteJSON(w, http.StatusOK, post)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var postResp services.Post
	err := json.NewDecoder(r.Body).Decode(&postResp)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	// TODO: Handle Error
	_ = helpers.WriteJSON(w, http.StatusOK, postResp)
	postCreated, err := models.Post.CreatePost(postResp)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		// TODO: Handle Error
		_ = helpers.WriteJSON(w, http.StatusOK, postCreated)
	}
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	id := chi.URLParam(r, "id")
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, post)
	postObj, err := post.UpdatePost(id, post)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusOK, postObj)
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := models.Post.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_ = helpers.WriteJSON(w, http.StatusOK, "Deleted!")
}
