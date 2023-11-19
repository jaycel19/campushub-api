package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	//"github.com/aws/aws-sdk-go/aws"

	//"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/services"
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
		limit = 20
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
	const maxRequestSize = 10 * 1024 * 1024 // 10mb
	var postReq services.PostRequest

	err := json.NewDecoder(r.Body).Decode(&postReq)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	// TODO: UPLOAD IMG TO AWS3
	// imageData := postReq.ImageData
	// filename, err := uuid.NewRandom()
	// if err != nil {
	// 	helpers.MessageLogs.ErrorLog.Println(err)
	// }
	// TODO: ERROR access denied
	// err = util.UploadImageToS3(bytes.NewReader(imageData), filename.String())

	// if err != nil {
	// 	fmt.Println(err)
	// 	helpers.MessageLogs.ErrorLog.Println(err)
	// 	// Handle error and send an appropriate response
	// 	_ = helpers.WriteJSON(w, http.StatusInternalServerError, "Failed to upload image")
	// 	return
	// }

	// imageUrl := "https://" + os.Getenv("AWS_S3_BUCKET") + ".s3.amazonaws.com/" + filename.String()
	imageUrl := "placeholder"
	postPayload := services.Post{
		Author:      postReq.Author,
		Image:       imageUrl,
		PostContent: postReq.PostContent,
	}
	// TODO: Handle Error
	_ = helpers.WriteJSON(w, http.StatusOK, postPayload)
	postCreated, err := models.Post.CreatePost(postPayload)
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
