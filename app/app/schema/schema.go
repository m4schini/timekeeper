package schema

import (
	"reflect"
	"time"

	"github.com/zitadel/schema"
)

func NewDecoder(options ...DecoderOption) *schema.Decoder {
	var decoder = schema.NewDecoder()
	for _, apply := range options {
		apply(decoder)
	}

	return decoder
}

type DecoderOption func(decoder *schema.Decoder)

func WithTime(dateLayout string) func(decoder *schema.Decoder) {
	return func(decoder *schema.Decoder) {
		decoder.RegisterConverter(time.Now(), func(s string) reflect.Value {
			t, err := time.Parse(dateLayout, s)
			if err != nil {
				return reflect.Value{}
			}
			return reflect.ValueOf(t)
		})
	}
}
