package pages

import (
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"time"
	"timekeeper/app/auth"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"
)

func LoginPage() Node {
	return Shell("",
		components.PageHeader(model.EventModel{}),
		Main(Style("width: 100%; height: 100%; display: flex; justify-content: center; align-items: center"),
			Form(Method("POST"), Action("/_/login"), Class("form"),
				Input(Type("text"), Name("username"), Placeholder("username")),
				Input(Type("password"), Name("password"), Placeholder("password")),
				Input(Type("submit"), Value("Login")),
			),
		),
	)
}

type LoginPageRoute struct {
	Auth auth.Authenticator
}

func (l *LoginPageRoute) Method() string {
	return http.MethodGet
}

func (l *LoginPageRoute) Pattern() string {
	return "/login"
}

func (l *LoginPageRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	page := LoginPage()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		render.HTML(log, writer, request, page)
	})
}

type LoginRoute struct {
	Auth auth.Authenticator
}

func (l *LoginRoute) Method() string {
	return http.MethodPost
}

func (l *LoginRoute) Pattern() string {
	return "/login"
}

func (l *LoginRoute) Handler() http.Handler {
	log := components.Logger(l)
	authy := l.Auth
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		var (
			username = request.PostFormValue("username")
			password = request.PostFormValue("password")
		)

		log.Debug("authenticating user")
		token, err := authy.AuthenticateUser(username, password)
		if err != nil {
			log.Warn("failed login", zap.Error(err))
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
		}

		log.Debug("user authenticated. setting cookie")
		http.SetCookie(writer, &http.Cookie{
			Name:     "SESSION",
			Value:    token,
			Path:     "/",
			Domain:   "",
			Expires:  time.Now().Add(18 * time.Hour),
			MaxAge:   int((18 * time.Hour).Seconds()),
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	})
}

type LogoutRoute struct {
}

func (l *LogoutRoute) Method() string {
	return http.MethodGet
}

func (l *LogoutRoute) Pattern() string {
	return "/logout"
}

func (l *LogoutRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "failed to parse form", err)
			return
		}

		log.Debug("user authenticated. setting cookie")
		http.SetCookie(writer, &http.Cookie{
			Name:     "SESSION",
			Value:    "",
			Path:     "/",
			Domain:   "",
			Expires:  time.Now(),
			MaxAge:   0,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	})
}
