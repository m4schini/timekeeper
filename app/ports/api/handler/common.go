package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"raumzeitalpaka/app/database/model"
	"strings"

	"go.uber.org/zap"
)

func InternalServerErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func Decode[T any](r *http.Request) (request T, err error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return request, fmt.Errorf("invalid content-type")
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

func Encode(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		zap.L().Error("failed to encode response", zap.Error(err))
	}
}

func ParseRolesQuery(query url.Values, userIsOrganizer bool) (roles []model.Role, hasRoles bool) {
	hasRoles = query.Has("role")
	rolesStrs := strings.Split(query.Get("role"), ",")

	if !hasRoles {
		if userIsOrganizer {
			rolesStrs = []string{string(model.RoleOrganizer), string(model.RoleMentor), string(model.RoleParticipant)}
		} else {
			rolesStrs = []string{string(model.RoleParticipant)}
		}
	}

	roles = make([]model.Role, len(rolesStrs))
	for i, role := range rolesStrs {
		roles[i] = model.RoleFrom(role)
	}

	return roles, hasRoles
}
