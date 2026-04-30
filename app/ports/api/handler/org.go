package handler

import (
	"net/http"
	"raumzeitalpaka/app/database/command"
)

func CreateOrg(create command.CreateOrganisation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func GetOrg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func UpdateOrg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := requireAuthentication(ctx)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func DeleteOrg() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

	}
}
