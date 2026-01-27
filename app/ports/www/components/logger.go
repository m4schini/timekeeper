package components

import (
	"fmt"

	"go.uber.org/zap"
)

type route interface {
	Method() string
	Pattern() string
}

func Logger(route route) *zap.Logger {
	return zap.L().Named("www").With(zap.String("route", fmt.Sprintf("%v %v", route.Method(), route.Pattern())))
}
