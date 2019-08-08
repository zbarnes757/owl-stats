package app

import (
	"context"
	"net/http"
	"os"
	"owl-stats/models"
	u "owl-stats/utils"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type contextKey string

var contextKeyUserID = contextKey("user")

// JwtAuthentication is a middleware for validating JWTs
func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/v1/user/new", "/api/v1/user/login"}
		requestPath := r.URL.Path

		// Check if request does not need authentication,
		// serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		// Token is missing, returns with error code 403 Unauthorized
		if tokenHeader == "" {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		// The token normally comes in format `Bearer {token-body}`, we check if
		// the retrieved token matched this requirement
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		// Grab the token part, what we are truly interested in
		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		// Malformed token, returns with http code 403 as usual
		if err != nil {
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		// Token is invalid, maybe not signed on this server
		if !token.Valid {
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		// Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		// Useful for monitoring
		ctx := context.WithValue(r.Context(), contextKeyUserID, tk.UserID)
		r = r.WithContext(ctx)

		// proceed in the middleware chain!
		next.ServeHTTP(w, r)

	})
}

// CurrentUser will return the currently signed in user from the request
func CurrentUser(r *http.Request) uint {
	return r.Context().Value(contextKeyUserID).(uint)
}
