package local

import (
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func NewHandler(authy Authenticator) (r chi.Router, err error) {
	log := zap.L()
	var rateLimiter *rate.Limiter
	loginPage := LoginPage()
	rateLimiter = rate.NewLimiter(rate.Limit(1), 1)

	mux := chi.NewMux()
	mux.Get("/login", func(writer http.ResponseWriter, request *http.Request) {
		render.HTML(log, writer, request, loginPage)
	})
	mux.Post("/login", func(writer http.ResponseWriter, request *http.Request) {
		err := rateLimiter.Wait(request.Context())
		if err != nil {
			log.Info("login was cancelled while waiting on rate limiter")
			return
		}
		err = request.ParseForm()
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
			return
		}

		log.Debug("user authenticated. setting cookie")
		auth.SetSessionCookie(writer, token)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	})
	mux.HandleFunc("/logout", func(writer http.ResponseWriter, request *http.Request) {
		log.Debug("clearing SESSSION Cookie")
		auth.ClearSessionCookie(writer)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		//writer.Write([]byte("logged out"))
	})
	return mux, nil
}

func LoginPage() Node {
	return components.Shell("",
		components.PageHeader(model.EventModel{}),
		Main(Style("width: 100%; height: 100%; display: flex; justify-content: center; align-items: center"),
			Form(Method("POST"), Action("/login"), Class("form"),
				Input(Type("text"), Name("username"), Placeholder("username")),
				Input(Type("password"), Name("password"), Placeholder("password")),
				Input(Type("submit"), Value("Login")),
			),
		),
	)
}
