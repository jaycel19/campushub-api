package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/services"
	"github.com/jaycel19/campushub-api/util"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var userCreds services.UserCreds
	err := json.NewDecoder(r.Body).Decode(&userCreds)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	user, err := models.User.UserLogin(userCreds.Username)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	err = util.CheckPassword(userCreds.Password, user.Password)
	if err != nil {
		_ = helpers.WriteJSON(w, http.StatusForbidden, helpers.Envelope{"Error": "Password not match!"})
		return
	}

	token, tokenId, err := util.GenerateAccessToken(user.Username)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	refreshToken, err := util.GenerateRefreshToken(token, tokenId)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	// New Session
	sessionPayload := services.Session{
		ID:           tokenId,
		Username:     user.Username,
		RefreshToken: refreshToken,
	}
	session, err := models.Session.CreateSession(sessionPayload)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{
		// TODO: CREATE A struct for loginresponse
		"session_id":     session.ID,
		"token":          token,
		"user":           user,
		"RefreshToken":   refreshToken,
		"TokenExpiresAt": session.ExpiresAt,
	},
	)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var user services.User
	all, err := user.GetAllUsers()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}

	_ = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"users": all})
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := models.User.GetUserById(username)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userResp services.User
	err := json.NewDecoder(r.Body).Decode(&userResp)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, userResp)
	userCreated, err := models.User.CreateUser(userResp)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		_ = helpers.WriteJSON(w, http.StatusOK, userCreated)
	}
}
