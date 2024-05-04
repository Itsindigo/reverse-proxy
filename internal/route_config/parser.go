package route_config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type HttpMethod string

const (
	get     HttpMethod = "GET"
	put     HttpMethod = "PUT"
	patch   HttpMethod = "PATCH"
	post    HttpMethod = "POST"
	delete  HttpMethod = "DELETE"
	options HttpMethod = "OPTIONS"
)

type RateLimitStrategy string

const (
	token_bucket = "token_bucket"
)

type Target struct {
	Host   string
	Port   string
	Path   string
	Method HttpMethod
}

type RateLimit struct {
	RateLimitStrategy RateLimitStrategy
}

type Route struct {
	Path      string     `yaml:"path"`
	Method    HttpMethod `yaml:"method"`
	Target    Target     `yaml:"target"`
	RateLimit RateLimit  `yaml:"rate_limit"`
}

type RouteConfig struct {
	version string  `yaml:"version"`
	Routes  []Route `yaml:"routes"`
}

type RouteMap struct {
	Routes []Route `yaml:"routes"`
}

func readConfigFile(path string) ([]Route, error) {
	var routeConfig RouteConfig
	data, err := ioutil.ReadFile(path)

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
