package pages

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"
	"strconv"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UserAccountPage(user model.UserModel, memberships []model.UserOrganisationMembership) Node {
	return components.Shell("",
		components.PageHeader(model.EventModel{}),
		Main(
			H1(Textf("User")),
			Div(
				Ul(
					Li(Textf("ID: %v", user.ID)),
					Li(Textf("Display Name: %v", user.DisplayName)),
					Li(Textf("Login Name: %v", user.LoginName)),
					Li(Textf("Last Login: %v", user.LastLogin)),
				),
			),
			H1(Textf("Groups")),
			Div(
				Ul(
					Map(memberships, func(membership model.UserOrganisationMembership) Node {
						return Li(Textf("%v: %v", membership.Name, membership.Role))
					}),
				),
			),
		),
	)
}

type UserAccountPageRoute struct {
	GetUser        query.GetUser
	GetMemberships query.GetUserOrganisations
}

func (l *UserAccountPageRoute) Method() string {
	return http.MethodGet
}

func (l *UserAccountPageRoute) Pattern() string {
	return "/user/{user}"
}

func (l *UserAccountPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			ctx         = request.Context()
			userParam   = chi.URLParam(request, "user")
			userID, err = strconv.ParseInt(userParam, 10, 64)
		)
		if err != nil {
			render.Error(log, writer, http.StatusBadRequest, "invalid user id", err)
			return
		}

		user, err := l.GetUser.Query(ctx, query.GetUserRequest{ID: int(userID)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get user", err)
			return
		}

		memberships, err := l.GetMemberships.Query(ctx, query.GetUserOrganisationsRequest{UserID: int(userID)})
		if err != nil {
			render.Error(log, writer, http.StatusInternalServerError, "failed to get memberships", err)
			return
		}

		page := UserAccountPage(user, memberships)

		render.HTML(log, writer, request, page)
	})
}

type MeAccountPageRoute struct{}

func (l *MeAccountPageRoute) Method() string {
	return http.MethodGet
}

func (l *MeAccountPageRoute) Pattern() string {
	return "/me"
}

func (l *MeAccountPageRoute) Handler() http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userId, _ := auth.UserFrom(request)

		http.Redirect(writer, request, fmt.Sprintf("/user/%v", userId), http.StatusFound)
	})
}
