package schema

import "net/http"

type Decoder interface {
	Decode(dst interface{}, src map[string][]string) error
}

func ParseForm[T any](decoder Decoder, r *http.Request) (form T, err error) {
	err = r.ParseForm()
	if err != nil {
		return form, err
	}

	err = decoder.Decode(&form, r.Form)
	return form, err
}
