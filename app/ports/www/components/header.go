package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/config"
	"time"

	. "maragu.dev/gomponents"
	htmx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func PageHeader(event model.EventModel) Node {
	now := time.Now()
	time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), config.Timezone())

	return Header(Class("page-header"),
		Logo(event.Name, event.ID),
		UserWidget(),
	)
}

func Logo(name string, eventId int) Node {
	if name == "" {
		name = "Raumzeitalpaka"
	}
	href := "/"
	if eventId != 0 {
		href = fmt.Sprintf("/event/%v", eventId)
	}
	return H1(Class("logo"),
		A(Text(name), Href(href)),
	)
}

func UserWidget() Node {
	return Div(Class("last-change"), htmx.Get(`/_/user/widget`), htmx.Trigger("load"),
		Div(Text("loading")),
	)
}

func UserWidgetLoggedOut() Node {
	return Div(Class("last-change"),
		A(Href("/login"), Text("Login")),
	)
}

func UserWidgetLoggedIn(user model.UserModel) Node {
	return Div(Class("last-change"),
		P(Text("Signed in as "), Strong(Text(user.LoginName))),
		A(Href("/logout"), Text("Logout")),
	)
}

type UserHeaderWidgetRoute struct {
	User query.GetUser
}

func (l *UserHeaderWidgetRoute) Method() string {
	return http.MethodGet
}

func (l *UserHeaderWidgetRoute) Pattern() string {
	return "/user/widget"
}

func (l *UserHeaderWidgetRoute) Handler() http.Handler {
	loggedOut := UserWidgetLoggedOut()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userId, _ := auth.UserFrom(request)
		user, err := l.User.Query(query.GetUserRequest{ID: userId})
		if err != nil {
			loggedOut.Render(writer)
			return
		}

		UserWidgetLoggedIn(user).Render(writer)
	})
}
