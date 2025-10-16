package components

import (
	"fmt"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/middleware"
	"timekeeper/ports/www/render"
)

func UserForm() Node {
	return Form(Method("POST"), Action("/_/user/new"), Class("form"),
		Div(
			Label(For("username"), Text("Username")),
			Input(Type("text"), Name("username"), Placeholder("Username"), Required()),
		),

		Div(
			Label(For("password"), Text("Password")),
			Input(Type("password"), Name("password"), Placeholder("Password"), Required()),
		),

		Input(Type("submit"), Value("Create")),
	)
}

type CreateUserRoute struct {
	DB *database.Database
}

func (l *CreateUserRoute) Method() string {
	return http.MethodPost
}

func (l *CreateUserRoute) Pattern() string {
	return "/user/new"
}

func (l *CreateUserRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	commands := l.DB.Commands
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.IsOrganizer(request) {
			render.RenderError(log, writer, http.StatusUnauthorized, "unauthorized request detected", nil)
			return
		}

		err := request.ParseForm()
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			usernameParam = request.PostFormValue("username")
			passwordParam = request.PostFormValue("password")
		)
		model, err := ParseCreateUserModel(usernameParam, passwordParam)
		if err != nil {
			render.RenderError(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}
		log.Debug("parsed create user form", zap.Any("model", model))

		id, err := commands.CreateUser(model)
		if err != nil {
			render.RenderError(log, writer, http.StatusInternalServerError, "failed to create user", err)
			return
		}
		log.Debug("created user", zap.Int("id", id))

		http.Redirect(writer, request, fmt.Sprintf("/event/%v", id), http.StatusSeeOther)
	})
}

func ParseCreateUserModel(username, password string) (model.CreateUserModel, error) {
	return model.CreateUserModel{
		LoginName:    username,
		PasswordHash: password,
	}, nil
}
