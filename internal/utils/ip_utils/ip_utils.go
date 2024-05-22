package ip_utils

import (
	"fmt"
	"net"
	"reflect"
)

func IsLocalhost(ip net.IP) bool {
	return reflect.DeepEqual(ip, net.ParseIP("::1")) || reflect.DeepEqual(ip, net.ParseIP("0:0:0:0:0:0:0:1")) || reflect.DeepEqual(ip, net.ParseIP("127.0.0.1")) || reflect.DeepEqual(ip, net.ParseIP("0.0.0.0"))
}

func GetIP(remoteAddr string, forwardHeader string) (string, error) {
	var host string
	ip, _, err := net.SplitHostPort(remoteAddr)

	if err != nil {
		return "", fmt.Errorf("userip: %q is not IP:port", remoteAddr)
	}

	userIP := net.ParseIP(ip)

	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", remoteAddr)
	}

	if IsLocalhost(userIP) {
		host = "localhost"
	} else {
		host = userIP.String()
	}

	/* Forward header takes precedence if defined  */
	if len(forwardHeader) > 0 {
		host = forwardHeader
	}

	return host, nil
}
