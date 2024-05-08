package proxy_configuration

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"

	"github.com/itsindigo/reverse-proxy/internal/constants"
)

type Target struct {
	Host   string
	Port   string
	Path   string
	Method constants.HttpMethod
}

type RateLimit struct {
	RequestsPerMinute int `yaml:"requests_per_minute"`
}

type Route struct {
	Path      string               `yaml:"path"`
	Method    constants.HttpMethod `yaml:"method"`
	Target    Target               `yaml:"target"`
	RateLimit RateLimit            `yaml:"rate_limit"`
}

type RouteConfig struct {
	Routes          []Route   `yaml:"routes"`
	GlobalRateLimit RateLimit `yaml:"global_rate_limit"`
}

type RouteMap struct {
	Routes []Route `yaml:"routes"`
}

func readConfigFile(path string) ([]Route, error) {
	var routeConfig RouteConfig
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, &RouteMapConfigError{Err: err}
	}

	yaml.Unmarshal(data, &routeConfig)

	return routeConfig.Routes, nil
}

func Load(config_file_path string) ([]Route, error) {
	routes, err := readConfigFile(config_file_path)

	if err != nil {
		log.Fatalf("Error reading route map config file: %v", err)
		return nil, err
	}

	return routes, nil
}
