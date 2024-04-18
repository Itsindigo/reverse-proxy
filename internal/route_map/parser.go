package route_map_parser

import "fmt"

type HttpMethod string

const (
	Get     HttpMethod = "GET"
	Put     HttpMethod = "PUT"
	Patch   HttpMethod = "PATCH"
	Post    HttpMethod = "POST"
	Delete  HttpMethod = "DELETE"
	Options HttpMethod = "OPTIONS"
)

type Target struct {
	host   string
	port   string
	path   string
	method HttpMethod
}

type Route struct {
	Path   string `yaml:"path"`
	Target Target `yaml:"target"`
}

type RouteMap struct {
	Routes []Route `yaml:"routes"`
}

func Parse(config_file_path string) {
	fmt.Printf("Received Config File Path: %s\n", config_file_path)
}
