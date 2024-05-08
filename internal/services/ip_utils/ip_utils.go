package ip_utils

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
)

func isLocalhost(ip net.IP) bool {
	return reflect.DeepEqual(ip, net.ParseIP("::1")) || reflect.DeepEqual(ip, net.ParseIP("127.0.0.1")) || reflect.DeepEqual(ip, net.ParseIP("0.0.0.0"))
}

func GetIpRequestKey(r *http.Request) (string, error) {
	var hostKey string
	ip, port, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}

	userIP := net.ParseIP(ip)

	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}

	if isLocalhost(userIP) {
		hostKey = "localhost"
	} else {
		hostKey = string(userIP)
	}

	/* Forward header takes precedence if defined  */
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		hostKey = forward
	}

	return fmt.Sprintf("%s:%s", hostKey, port), nil
}
