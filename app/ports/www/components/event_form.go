package components

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/schema"
	"regexp"
	"time"

	"codeberg.org/aur0ra/form"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const (
	EventSlugPattern = `(?i)^[a-z0-9-]+$` // (?i) makes it case-insensitive
	EventDateLayout  = "02.01.2006"
)

var (
	eventSlugRegex   = regexp.MustCompile(EventSlugPattern)
	eventFormDecoder = form.MustNewDecoder[EventFormSchema]()
)

func EventUpdateForm(event model.EventModel) Node {
	return eventForm(&event, "POST", fmt.Sprintf("/event/%v/edit", event.ID), "Update")
}

func EventCreateForm() Node {
	return eventForm(nil, "POST", "/event/new", "Create")
}

type EventFormSchema struct {
	Event    int    `form:"event"`
	Name     string `form:"name,required"`
	Slug     string `form:"slug,required"`
	Start    Date   `form:"start,required"`
	End      Date   `form:"end,required"`
	Setup    int    `form:"setup,required"`
	Teardown int    `form:"teardown,required"`
}

func (s EventFormSchema) Validate(requireEventId bool) error {
	if requireEventId && s.Event <= 0 {
		return schema.InvalidFieldValueErr("event")
	}

	if !eventSlugRegex.MatchString(s.Slug) {
		return schema.InvalidFieldValueErr("slug")
	}

	return nil
}

func eventForm(event *model.EventModel, method, action, actionText string) Node {
	hasModel := event != nil
	eventId := -1
	if hasModel {
		eventId = event.ID
	}
	return Form(Method(method), Action(action), Class("form"),
		Input(Type("hidden"), Name("event"), Value(fmt.Sprintf("%v", eventId))),

		Div(Class("param"),
			Label(For("name"), Text("Name")),
			Input(Type("text"), Name("name"), Placeholder("Jugend hackt 2042"), Required(), Iff(hasModel, func() Node {
				return Value(event.Name)
			})),
		),

		Div(Class("param"),
			Label(For("slug"), Text("URL Slug")),
			Input(Type("text"), Name("slug"), Placeholder("jh42"), Pattern("^[A-Za-z0-9-]+$"), Required(), Iff(hasModel, func() Node {
				return Value(event.Slug)
			})),
		),

		Div(Class("param"),
			Label(For("start"), Text("Erster Tag")),
			Input(Type("text"), Name("start"),
				Placeholder("13.12.2042"),
				Pattern(`^(0[1-9]|[12][0-9]|3[01])\.(0[1-9]|1[0-2])\.(\d{4})$`),
				Required(), Iff(hasModel, func() Node {
					return Value(event.Start.Format("02.01.2006"))
				})),
		),

		Div(Class("param"),
			Label(For("end"), Text("Letzter Tag")),
			Input(Type("text"), Name("end"),
				Placeholder("42.12.2161"),
				Pattern(`^(0[1-9]|[12][0-9]|3[01])\.(0[1-9]|1[0-2])\.(\d{4})$`),
				Required(), Iff(hasModel, func() Node {
					return Value(event.End.Format("02.01.2006"))
				})),
		),

		Div(Class("param"),
			Label(For("setup"), Text("Aufbau (0 = Aufbau am ersten Tag, 1 = Ein Tag vor dem Event, ...)")),
			Input(Type("text"), Name("setup"),
				Placeholder("0"),
				Required(), Iff(hasModel, func() Node {
					return Value(fmt.Sprintf("%v", event.Setup))
				})),
		),

		Div(Class("param"),
			Label(For("teardown"), Text("Abbau (0 = Abbau am letzten Tag, 1 = Ein Tag nach dem Event, ...)")),
			Input(Type("text"), Name("teardown"),
				Placeholder("0"),
				Required(), Iff(hasModel, func() Node {
					return Value(fmt.Sprintf("%v", event.Teardown))
				})),
		),

		Input(Type("submit"), Value(actionText)),
	)
}

func DecodeEventForm(r *http.Request, requireEventId bool) (EventFormSchema, error) {
	err := r.ParseForm()
	if err != nil {
		return EventFormSchema{}, err
	}

	form, err := eventFormDecoder.Decode(r.Form)
	if err != nil {
		return EventFormSchema{}, err
	}

	return form, form.Validate(requireEventId)
}

type Date time.Time

func (Date) ParseForm(field string) (any, error) {
	t, err := time.Parse(EventDateLayout, field)
	return Date(t), err
}
