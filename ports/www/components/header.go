package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"time"
	"timekeeper/app/database/model"
)

func PageHeader(event model.EventModel, withActions bool) Node {
	//g := Group{}
	//for i := 0; i < event.TotalDays; i++ {
	//	g = append(g, A(Href(fmt.Sprintf("/event/%v/%v", event.ID, i)), Textf("Tag %v", i+1)))
	//}

	menu := Div(Class("menu"),
		//A(Href("/location"), Text("Karte")),
		A(Href(fmt.Sprintf("/event/%v?role=Organizer,Mentor,Participant", event.ID)), Text("Orga")),
		A(Href(fmt.Sprintf("/event/%v?role=Mentor,Participant", event.ID)), Text("Mentor*innen")),
		A(Href(fmt.Sprintf("/event/%v?role=Participant", event.ID)), Text("Teilnehmer*innen")),
	)

	return Header(Class("page-header"),
		Logo(event.Name, event.ID),
		If(withActions, menu),
		Div(Class("last-change"),
			Text("Generated:"),
			Br(),
			Text(time.Now().Format(time.RFC822)),
		),
	)
}
