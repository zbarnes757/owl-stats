package controllers

import (
	"net/http"
	"owl-stats/models"

	"github.com/google/jsonapi"
)

// CreateAccount will sign up a new user
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := new(models.Account)

	if err := jsonapi.UnmarshalPayload(r.Body, account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	account.Create()

	w.Header().Set("Content-Type", jsonapi.MediaType)
	// w.WriteHeader(http.StatusCreated)
	if err := jsonapi.MarshalPayload(w, account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
