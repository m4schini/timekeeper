package schema

import "fmt"

type InvalidFieldValueErr string

func (e InvalidFieldValueErr) Error() string {
	return fmt.Sprintf("invalid field value: %v", e)
}
