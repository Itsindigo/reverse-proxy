package route_config

import (
	"errors"
	"fmt"
)

var ErrReadingRoutesConfig = errors.New("error reading route map config file")

type RouteMapConfigError struct {
	Err error
}

func (e *RouteMapConfigError) Error() string {
	return fmt.Sprintf("Error reading route map config: %v", e.Err)
}
