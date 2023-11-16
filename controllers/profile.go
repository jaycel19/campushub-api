package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/services"
)

type Background struct {
	BgColor string `json:"profile_background"`
}

// GET/posts
func GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	var profile services.Profile
	all, err := profile.GetAllProfiles()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"profiles": all})
}

func GetProfileByUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	profile, err := models.Profile.GetProfileByUser(username)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	// TODO: handle error
	_ = helpers.WriteJSON(w, http.StatusOK, profile)
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profileReq services.ProfileRequest

	err := json.NewDecoder(r.Body).Decode(&profileReq)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	profilePayload := services.Profile{
		Name:     profileReq.Name,
		Username: profileReq.Username,
		Age:      profileReq.Age,
		Program:  profileReq.Program,
		Year:     profileReq.Year,
	}
	// TODO: Handle Error
	profileCreated, err := models.Profile.CreateProfile(profilePayload)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
		// TODO: Handle Error
	}
	_ = helpers.WriteJSON(w, http.StatusOK, profileCreated)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profile services.Profile
	username := chi.URLParam(r, "username")
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	profileObj, err := profile.UpdateProfile(username, profile)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, profileObj)
}

func ProfileChangeBackground(w http.ResponseWriter, r *http.Request) {
	var profile services.Profile
	var background Background
	username := chi.URLParam(r, "username")
	err := json.NewDecoder(r.Body).Decode(&background)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = profile.ProfileChangeBackground(username, background.BgColor)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"success": "background changed!", "bg": background})
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	err := models.Profile.DeleteProfile(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = helpers.WriteJSON(w, http.StatusOK, "Deleted!")
}
