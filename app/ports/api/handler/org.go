package handler

import (
	"net/http"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func GetOrg(organisation query.GetOrganisation, orgBySlug query.GetOrganisationBySlug,
	members query.GetOrganisationMembers,
	events query.GetEventsByOrganisation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgParam := chi.URLParam(r, "org")
		var org model.OrganisationModel
		orgID, err := strconv.ParseInt(orgParam, 10, 64)
		if err != nil {
			org, err = orgBySlug.Query(ctx, query.GetOrganisationBySlugRequest{Slug: orgParam})
		} else {
			org, err = organisation.Query(ctx, query.GetOrganisationRequest{ID: int(orgID)})
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//TODO prevent number only org slugs

		//ms, err := members.Query(ctx, query.GetOrganisationMembersRequest{OrganisationID: org.ID})
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//
		//es, err := events.Query(ctx, query.GetEventsByOrganisationRequest{OrganisationID: org.ID})
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		Encode(w, org)
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
