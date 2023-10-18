package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/util"
)

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	RefreshToken   string  `json:"refresh_token"`
	TokenExpiresAt float64 `json:"token_expires_at"`
}

func RenewToken(w http.ResponseWriter, r *http.Request) {
	var tokenRequest TokenRequest
	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	refreshId, err := util.VerifyToken(tokenRequest.RefreshToken)
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"error": err})
		return
	}
	session, err := models.Session.GetSessionById(refreshId)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"error": "Server error"})
		return
	}

	if session.IsBlocked {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"error": "Token invalid"})
		return
	}

	token, _, err := util.GenerateAccessToken(session.Username)
	// TODO make exp a constant or a env variable
	exp := time.Now().Add(time.Minute * 15).Unix()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"token": token, "token_expires_at": exp})
}
