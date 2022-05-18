package yoorel

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/c0c0n3/resto/util/err"
)

const (
	DEFAULT_HTTP_PORT  = 80
	DEFAULT_HTTPS_PORT = 443
)

// Host and port part of a URL. The host can be either a host name or
// IP address.
type HostAndPort struct {
	h string
	p int
}

// Parse the given string as an integer between 0 and 65535.
func ParsePort(p string) (int, error) {
	p = strings.TrimSpace(p)
	if port, err := strconv.Atoi(p); err == nil {
		if 0 <= port && port <= 65535 {
			return port, nil
		}
	}
	return 0, err.Mk[InvalidPort]("%s", p)
}

var hostnameRx = regexp.MustCompile(
	`^(([a-zA-Z0-9_-]){1,63}\.)*([a-zA-Z0-9_-]){1,63}$`)

// This article explains quite well what makes up a valid hostname:
// - https://en.wikipedia.org/wiki/Hostname

// Is the given argument a valid host? i.e. something you can use as the
// host part of a URL.
func IsHostname(host string) error {
	if 0 < len(host) && len(host) < 254 {
		if net.ParseIP(host) != nil || hostnameRx.MatchString(host) {
			return nil
		}
	}
	return err.Mk[InvalidHostname]("%s", host)
}

// Is the input in the format "host:port" or does it only have the host
// part?
func HasPort(hp string) bool {
	if _, _, err := net.SplitHostPort(hp); err == nil {
		return true
	}
	return false
}

// Parse an input in the format "host:port" where host and port are the URL
// components we know and love.
func ParseHostAndPort(hp string) (*HostAndPort, error) {
	hp = strings.TrimSpace(hp) // (1)
	if host, portString, err := net.SplitHostPort(hp); err != nil {
		return nil, err
	} else {
		if err := IsHostname(host); err != nil { // (2)
			return nil, err
		}
		if port, err := ParsePort(portString); err != nil { // (3)
			return nil, err
		} else {
			return &HostAndPort{host, port}, nil
		}
	}

	// (1) SplitHostPort doesn't trim space, e.g.
	//       SplitHostPort(" h:1 ") == (" h", "1 ", nil)
	// (2) SplitHostPort doesn't check the host part is a valid IP4 or IP6 or
	//     a valid hostname e.g.
	//       SplitHostPort(":123") == ("", "123", nil)
	//       SplitHostPort("??:123") == ("??", "123", nil)
	// (3) SplitHostPort doesn't check the port range, e.g.
	//       SplitHostPort("h:123456789") == ("h", "123456789", nil)
}

// The host name or IP address.
func (d *HostAndPort) Host() string {
	return d.h
}

// The port.
func (d *HostAndPort) Port() int {
	return d.p
}

// Stringify host name and port using the format "host:port".
func (d *HostAndPort) String() string {
	return fmt.Sprintf("%s:%d", d.h, d.p)
}
