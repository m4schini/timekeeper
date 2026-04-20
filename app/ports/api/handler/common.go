package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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
