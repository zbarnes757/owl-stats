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

	w.Header().Set("Content-Type", jsonapi.MediaType)
	if err := account.Create(); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Validation Error",
			Detail: err.Error(),
			Status: "400",
		}})
		return
	}

	if err := jsonapi.MarshalPayload(w, account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
