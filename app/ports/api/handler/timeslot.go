package handler

import (
	"net/http"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/query"
)

func CreateTimeslot(createTimeslot command.CreateTimeslot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func GetTimeslot(getTimeslot query.GetTimeslot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func UpdateTimeslot(updateTimeslot command.UpdateTimeslot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}
