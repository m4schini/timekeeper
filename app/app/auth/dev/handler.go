package dev

import (
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/auth/local"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/render"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewHandler(upsertUser command.UpsertUser, authenticator local.Authenticator) (r chi.Router, err error) {
	log := zap.L()
	loginPage := local.LoginPage()

	mux := chi.NewMux()
	mux.Get("/login", func(writer http.ResponseWriter, request *http.Request) {
		render.HTML(log, writer, request, loginPage)
	})
	mux.Post("/login", func(writer http.ResponseWriter, request *http.Request) {
		err = request.ParseForm()
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			username = request.PostFormValue("username")
			password = request.PostFormValue("password")
		)

		pwHash, err := local.GeneratePasswordHash(password, &local.DefaultPasswordParams)
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to generate password", err)
			return
		}

		log.Debug("authenticating user")
		_, err = upsertUser.Execute(command.UpsertUserRequest{
			ID:           1,
			LoginName:    username,
			PasswordHash: pwHash,
			Role:         model.RoleOrganizer,
		})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to upsert", err)
			return
		}

		token, err := authenticator.AuthenticateUser(username, password)
		if err != nil {
			log.Warn("failed login", zap.Error(err))
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}

		log.Debug("user authenticated. setting cookie")
		auth.SetSessionCookie(writer, token)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	})
	mux.HandleFunc("/logout", func(writer http.ResponseWriter, request *http.Request) {
		log.Debug("clearing SESSSION Cookie")
		auth.ClearSessionCookie(writer)
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
	})
	return mux, nil
}
