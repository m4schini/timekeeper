package query

func (q *Queries) GetEventIdBySlug(slug string) (id int, err error) {
	row := q.DB.QueryRow(`SELECT id FROM raumzeitalpaka.events WHERE slug = $1`, slug)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
