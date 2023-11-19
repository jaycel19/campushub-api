package middlewares

import (
	"net/http"
	"strings"

	"github.com/jaycel19/campushub-api/helpers"
	"github.com/jaycel19/campushub-api/services"
	"github.com/jaycel19/campushub-api/util"
)

var models services.Models

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"Error": "Missing Authorization Header."})
			return
		}
		splitToken := strings.Split(tokenHeader, "Bearer ")
		if len(splitToken) != 2 {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"Error": "Invalid Token Format."})
			return
		}
		tokenString := splitToken[1]
		tokenID, err := util.VerifyToken(tokenString)
		if err != nil {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"Message": "Invalid Token Format", "Error": err})
			return
		}
		session, err := models.Session.GetSessionById(tokenID)
		if err != nil {
			helpers.MessageLogs.ErrorLog.Println(err)
			helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"error": "Internal Server error"})
			return
		}

		if session.IsBlocked {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"error": "Token invalid"})
			return
		}
		user, err := models.User.GetUserById(session.Username)
		if err != nil {
			helpers.MessageLogs.ErrorLog.Println(err)
			return
		}

		if user == nil {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"message": "User not found!"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
