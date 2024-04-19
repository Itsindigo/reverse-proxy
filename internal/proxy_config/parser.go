package proxy_config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type httpMethod string

const (
	get     httpMethod = "GET"
	put     httpMethod = "PUT"
	patch   httpMethod = "PATCH"
	post    httpMethod = "POST"
	delete  httpMethod = "DELETE"
	options httpMethod = "OPTIONS"
)

type Target struct {
	Host   string
	Port   string
	Path   string
	Method httpMethod
}

type Route struct {
	Path   string `yaml:"path"`
	Target Target `yaml:"target"`
}

type RouteConfig struct {
	version string  `yaml:"version"`
	Routes  []Route `yaml:"routes"`
}

type RouteMap struct {
	Routes []Route `yaml:"routes"`
}

func read_config_file(path string) ([]Route, error) {
	var routeConfig RouteConfig
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, &RouteMapConfigError{Err: err}
	}

	yaml.Unmarshal(data, &routeConfig)

	return routeConfig.Routes, nil
}

func Parse(config_file_path string) ([]Route, error) {
	routes, err := read_config_file(config_file_path)

	if err != nil {
		log.Fatalf("Error reading route map config file: %v", err)
		return nil, err
	}

	return routes, nil
}
