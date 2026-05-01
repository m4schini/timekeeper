//go:generate go tool oapi-codegen -config config.yaml ../../raumzeitalpaka.openapi.yaml
package api

import (
	"net/http"
)

func New(si ServerInterface, options ChiServerOptions) http.Handler {
	return HandlerWithOptions(si, options)
}
