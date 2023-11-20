package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/services"
	"github.com/jaycel19/campushub-api/util"
)

var models services.Models

// GET/posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	var post services.Post

	// Get the limit, offset, and topPost parameters from the query
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	topPostParam := r.URL.Query().Get("topPost")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	all, err := post.GetPosts(limit, offset, topPostParam)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		// Handle the error and return an appropriate response
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"error": "Internal Server Error"})
		return
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
	var postReq services.PostRequest
	// Extract JSON data from the form
	jsonBodyStr := r.FormValue("jsonBody")
	if jsonBodyStr == "" {
		helpers.MessageLogs.ErrorLog.Println("JSON data not found in the form")
		_ = helpers.WriteJSON(w, http.StatusBadRequest, "JSON data not found in the form")
		return
	}

	// Decode JSON data
	err := json.Unmarshal([]byte(jsonBodyStr), &postReq)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusBadRequest, "Failed to decode JSON data")
		return
	}

	// Parse multipart form data
	err = r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusBadRequest, "Failed to retrieve image from form data")
		return
	}
	defer file.Close()
	fileExt := filepath.Ext(header.Filename)
	// Generate a unique filename using UUID
	fileKey, err := uuid.NewRandom()
	filename := fileKey.String() + fileExt
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, "Failed to generate filename")
		return
	}

	// Upload image to AWS S3
	err = util.UploadImageToS3(file, filename)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, "Failed to upload image to S3")
		return
	}

	imageUrl := "https://campushub-beta.s3.amazonaws.com/" + filename
	postPayload := services.Post{
		Author:      postReq.Author,
		Image:       imageUrl,
		PostContent: postReq.PostContent,
	}

	// Create post in the database
	postCreated, err := models.Post.CreatePost(postPayload)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, "Failed to create post in the database")
		return
	}

	// Send a successful response
	_ = helpers.WriteJSON(w, http.StatusOK, postCreated)
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

func PostLike(w http.ResponseWriter, r *http.Request) {
	var post services.Post
	id := chi.URLParam(r, "id")
	err := post.PostLike(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	fmt.Println(id)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := models.Post.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_ = helpers.WriteJSON(w, http.StatusOK, "Deleted!")
}
