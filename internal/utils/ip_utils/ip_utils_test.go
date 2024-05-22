package ip_utils

import (
	"net"
	"testing"
)

func TestIsLocalhost(t *testing.T) {
	testCases := []struct {
		name string
		ip   net.IP
		want bool
	}{
		{"Should be true for IPV6 longhand", net.ParseIP("0:0:0:0:0:0:0:1"), true},
		{"Should be true for IPV6 shorthand", net.ParseIP("::1"), true},
		{"Should be true for 0.0.0.0", net.ParseIP("0.0.0.0"), true},
		{"Should be true for 127.0.0.1", net.ParseIP("127.0.0.1"), true},
		{"Should be false for localhost", net.ParseIP("localhost"), false},
		{"Should be false for 172.217.0.0", net.ParseIP("172.217.0.0"), false},
		{"Should be false for 2001:569:7bfa:3400:9cf4:eb79:10c5:9ac8", net.ParseIP("2001:569:7bfa:3400:9cf4:eb79:10c5:9ac8"), false},
		{"Should be false for 75.154.251.224", net.ParseIP("75.154.251.224"), false},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := IsLocalhost(test.ip)
			want := test.want

			if got != want {
				t.Errorf("Expected %t, got %t", want, got)
			}
		})
	}
}

func TestGetIP(t *testing.T) {
	testCases := []struct {
		name          string
		remoteAddr    string
		forwardHeader string
		want          string
	}{
		{"should extract forward header when present", "[::1]:54865", "75.154.251.224", "75.154.251.224"},
		{"should return IP when valid", "192.168.1.64:54909", "", "192.168.1.64"},
		{"should return localhost when called from local address", "[::1]:54865", "", "localhost"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, err := GetIP(test.remoteAddr, test.forwardHeader)
			want := test.want

			if err != nil {
				t.Errorf("Got error calling GetIP: %s", err)
			}

			if got != want {
				t.Errorf("Expected %s, got %s", want, got)
			}
		})
	}
}
