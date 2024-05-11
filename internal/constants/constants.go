package constants

type HttpMethod string

const (
	get     HttpMethod = "GET"
	put     HttpMethod = "PUT"
	patch   HttpMethod = "PATCH"
	post    HttpMethod = "POST"
	delete  HttpMethod = "DELETE"
	options HttpMethod = "OPTIONS"
)

func (m HttpMethod) String() string {
	return string(m)
}

type RedisKeyNamespaces string

/* Redis Key Namespaces */
const (
	UserHttpRequestLimit RedisKeyNamespaces = "UserHttpRequestLimit"
)
