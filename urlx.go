// Package urlx parses and normalizes URLs. It favors absolute paths.
package urlx

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/purell"
)

var (
	ErrInvalidURL       = errors.New("invalid URL")
	ErrInvalidURLHost   = errors.New("invalid hostname")
	ErrInvalidURLScheme = errors.New("invalid URL scheme")
)

func NormalizeString(s string) (string, error) {
	// Hack for "localhost", as net/url can't handle it.
	if len(s) == 9 && s == "localhost" ||
		len(s) > 9 && (s[0:10] == "localhost:" || s[0:10] == "localhost/") {
		s = "http://" + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	if u.Host == "" {
		// url.Parse("example.com") returns Host="", Path="example.com", eh..
		i := strings.IndexByte(u.Path, '/')
		if i == -1 {
			i = len(u.Path)
		}
		if i < 3 {
			// Path too short to contain the actual Host
			return s, ErrInvalidURLHost
		}
		// Move the actual host from Path to Host field
		u.Host = u.Path[:i]
		u.Path = u.Path[i:]
	}
	switch u.Scheme {
	case "http", "https", "HTTP", "HTTPS":
		// nop
	case "":
		// AnyScheme => http
		u.Scheme = "http"
	default:
		return s, ErrInvalidURLScheme
	}
	return purell.NormalizeURL(u, purell.FlagsUsuallySafeGreedy|purell.FlagRemoveDuplicateSlashes|purell.FlagLowercaseScheme|purell.FlagLowercaseHost), nil
}

func Validate(rawURL string, dnsCheck bool) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ErrInvalidURL
	}
	if parsed.Scheme == "" {
		return ErrInvalidURL
	}

	hostPort := strings.Split(parsed.Host, ":")
	hostname := hostPort[0]

	reg := regexp.MustCompile(`[\w\.\-]{2,64}\.[a-zA-Z0-9\-]{2,64}$`)

	if !reg.Match([]byte(hostname)) {
		return ErrInvalidURLHost
	}

	if dnsCheck {
		if _, err := net.ResolveIPAddr("ip", hostname); err != nil {
			return err
		}
	}

	return nil
}
