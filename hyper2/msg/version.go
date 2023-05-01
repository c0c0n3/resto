package msg

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/c0c0n3/resto/hyper2"
)

// HTTP protocol version.
// See:
// - https://www.rfc-editor.org/rfc/rfc9110.html#name-protocol-version

// An HTTP protocol version, e.g. HTTP/1.1.
type ProtocolVersion struct {
	major int
	minor int
}

// The major number component of the HTTP protocol version.
func (v ProtocolVersion) Major() int {
	return v.major
}

// The minor number component of the HTTP protocol version.
func (v ProtocolVersion) Minor() int {
	return v.minor
}

const httpVersionPrefix = "HTTP/"

func versionRegex() *regexp.Regexp {
	// (?i)^HTTP/(\d+)[.](\d+)$|^HTTP/(\d+)$
	expr := fmt.Sprintf(`(?i)^%s(\d+)[.](\d+)$|^%s(\d+)$`,
		httpVersionPrefix,
		httpVersionPrefix)
	return regexp.MustCompile(expr)
}

var httpVersionRegex = versionRegex()

// Parse the wire representation of the protocol version.
// Be liberal in what you accept. If there's no minor version assume
// it's 0. E.g. some servers output "HTTP/2" instead of "HTTP/2.0".
// Also, treat the "HTTP/" prefix as case insensitive.
//
// NOTE. "net/http" is much stricter. It won't parse a version without
// the minor number and won't parse a prefix that's not exactly "HTTP/".
// E.g.
//                                            M m ok
//     http.ParseHTTPVersion("HTTP/2.0")  ~~> 2 0 true
//     http.ParseHTTPVersion("Http/2.0")  ~~> 0 0 false
//     http.ParseHTTPVersion("HTTP/2")    ~~> 0 0 false
//
func ParseProtocolVersion(v string) (ProtocolVersion, error) {
	parsed := ProtocolVersion{0, 0}
	groups := httpVersionRegex.FindStringSubmatch(v)
	// "HTTP/1"     ~~> ["HTTP/1","","","1"]
	// "HTTP/21.32" ~~> ["HTTP/21.32","21","32",""]
	if len(groups) != 4 {
		return parsed, hyper2.UnparsableProtocolVersionErr(v)
	}

	major, minor, majorOnly := groups[1], groups[2], groups[3]
	n, _ := strconv.Atoi(major + majorOnly) // one of the two is empty
	parsed.major = n
	if len(minor) > 0 {
		n, _ := strconv.Atoi(minor)
		parsed.minor = n
	}

	return parsed, nil
}

// Format the HTTP protocol version as "HTTP/M.m" where M and m are,
// respectively the major and minor version numbers.
func (v ProtocolVersion) String() string {
	return fmt.Sprintf("%s%d.%d", httpVersionPrefix, v.major, v.minor)
}

// The "HTTP/1.1" version.
func HTTP_1_1() ProtocolVersion {
	return ProtocolVersion{major: 1, minor: 1}
}

// The "HTTP/2.0" version.
func HTTP_2_0() ProtocolVersion {
	return ProtocolVersion{major: 2, minor: 0}
}

// The "HTTP/3.0" version.
func HTTP_3_0() ProtocolVersion {
	return ProtocolVersion{major: 3, minor: 0}
}
