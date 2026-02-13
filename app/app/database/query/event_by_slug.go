package query

type GetEventBySlug Handler[GetEventBySlugRequest, int]

type GetEventBySlugRequest struct {
	Slug string
}

type GetEventBySlugHandler struct {
	DB Database
}

func (q *GetEventBySlugHandler) Query(request GetEventBySlugRequest) (id int, err error) {
	row := q.DB.QueryRow(`SELECT id FROM raumzeitalpaka.events WHERE slug = $1`, request.Slug)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
