package components

import (
	"fmt"
	"go.uber.org/zap"
	"timekeeper/ports/www"
)

func Logger(route www.Route) *zap.Logger {
	return zap.L().Named("www").With(zap.String("route", fmt.Sprintf("%v %v", route.Method(), route.Pattern())))
}
