package proxy_configuration

import (
	"fmt"
)

type RouteMapConfigError struct {
	Err error
}

func (e *RouteMapConfigError) Error() string {
	return fmt.Sprintf("Error reading route map config: %v", e.Err)
}
