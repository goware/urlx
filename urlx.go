// Package urlx parses and normalizes URLs. It can also resolve hostname to an IP address.
package urlx

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
)

var (
	ErrEmptyURLHost     = errors.New("empty hostname")
	ErrInvalidURLHost   = errors.New("invalid hostname")
	ErrUnresolvableHost = errors.New("unable to resolve hostname")
)

type URL struct {
	*url.URL
	Host string
	Port string
}

// Parse parses raw URL string into the URL struct. It uses net/url.Parse()
// internally, but it slightly changes it's behavior:
// 1. It forces the default scheme and port.
// 2. It favors absolute paths over relative ones, thus "example.com" is
//    parsed into url.Host instead of into url.Path.
// 3. It splits Host:Port into separate fields by default.
func Parse(raw string) (*URL, error) {
	var err error
	url := &URL{}

	// Force default http scheme, so net/url.Parse() doesn't
	// put both host and path into the (relative) path.
	if strings.Index(raw, "//") == 0 {
		// Leading double slashes (any scheme). Force http.
		raw = "http:" + raw
	}
	if strings.Index(raw, "://") == -1 {
		// Missing scheme. Force http.
		raw = "http://" + raw
	}

	// Use net/url.Parse() now.
	url.URL, err = url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if url.URL.Host == "" {
		return nil, ErrEmptyURLHost
	}
	url.Host = strings.ToLower(url.URL.Host)

	// Split Host:Port, if possible. Ignore anything inside [IPv6] brackets.
	// Don't use net.SplitHostPort, as it removes brackets from [IPv6] Host.
	if i := strings.LastIndex(url.Host, ":"); i != -1 && strings.Index(url.Host[i:], "]") == -1 {
		if len(url.Host) > i {
			url.Port = url.Host[i+1:]
		}
		url.Host = url.Host[:i]
	}

	// Force default port matching the http or https scheme.
	// TODO: Add hash table to support more schemes.
	if url.Port == "" {
		switch url.Scheme {
		case "http":
			url.Port = "80"
		case "https":
			url.Port = "443"
		}
	}

	return url, nil
}

// String returns URL struct as a string in human readable form:
// scheme://userinfo@host/path?query#fragment
// ^^^^^^^^^         ^^^^ (required parts)
func (url *URL) String() string {
	result := url.Scheme + "://"

	if url.User != nil {
		result += url.User.String() + "@"
	}

	// Show port only if it's not matching the scheme's default
	// port, eg. http://example.com:80 => http://example.com.
	// TODO: Add hash table to support more schemes.
	switch url.Scheme + ":" + url.Port {
	case "http:80", "https:443":
		result += url.Host
	default:
		result += url.Host + ":" + url.Port
	}

	if url.Path != "" {
		result += url.Path
	}

	if url.RawQuery != "" {
		result += "?" + url.RawQuery
	}

	if url.Fragment != "" {
		result += "#" + url.Fragment
	}

	return result
}

// Normalize returns normalized URL string.
// Behavior:
// 1. Remove unnecessary host dots.
// 2. Remove default port (http://localhost:80 becomes http://localhost).
// 3. Remove duplicate slashes.
// 4. Remove unnecessary dots from path.
// 5. Sort query parameters.
// 6. Decode host IP into decimal numbers.
// 7. Handle escape values.
func (url *URL) Normalize() (string, error) {
	var flags purell.NormalizationFlags
	flags |= purell.FlagRemoveUnnecessaryHostDots | purell.FlagRemoveDefaultPort
	flags |= purell.FlagDecodeDWORDHost | purell.FlagDecodeOctalHost | purell.FlagDecodeHexHost
	flags |= purell.FlagRemoveDuplicateSlashes | purell.FlagRemoveDotSegments
	flags |= purell.FlagUppercaseEscapes | purell.FlagDecodeUnnecessaryEscapes | purell.FlagEncodeNecessaryEscapes
	flags |= purell.FlagSortQuery

	normalized, err := purell.NormalizeURLString(url.String(), flags)
	if err != nil {
		return "", err
	}
	return normalized, nil
}

// Resolve resolves the URL host to its IP address.
func (url *URL) Resolve() (*net.IPAddr, error) {
	addr, err := net.ResolveIPAddr("ip", url.Host)
	if err != nil {
		return nil, ErrUnresolvableHost
	}
	return addr, nil
}
