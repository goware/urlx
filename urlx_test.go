package urlx_test

import (
	"testing"

	"github.com/goware/urlx"
)

var tests = []struct {
	in  string
	out string
	err bool
}{
	// Error out on missing host:
	{in: "", out: "", err: true},
	{in: "/", out: "/", err: true},
	{in: "/foo.html", out: "/foo.html", err: true},
	// Default scheme http://:
	{in: "localhost", out: "http://localhost"},
	{in: "localhost:8080", out: "http://localhost:8080"},
	{in: "example.com", out: "http://example.com"},
	{in: "1.example.com", out: "http://1.example.com"},
	{in: "subsub.sub.example.com", out: "http://subsub.sub.example.com"},
	// AnyScheme => http://:
	{in: "//example.com", out: "http://example.com"},
	// Supported schemes:
	{in: "http://example.com", out: "http://example.com"},
	{in: "https://example.com", out: "https://example.com"},
	{in: "HTTP://example.com", out: "http://example.com"},
	{in: "HTTPS://example.com", out: "https://example.com"},
	// Unsupported schemes:
	{in: "file://example.com", out: "file://example.com", err: true},
	{in: "mailto://user@example.com", out: "mailto://user@example.com", err: true},
	{in: "ftp://example.com", out: "ftp://example.com", err: true},
	{in: "git://example.com", out: "git://example.com", err: true},
	{in: "git+ssh://example.com", out: "git+ssh://example.com", err: true},
	{in: "ws://example.com", out: "ws://example.com", err: true},
	{in: "wss://example.com", out: "wss://example.com", err: true},
	{in: "magnet://example.com", out: "magnet://example.com", err: true},
	{in: "jabber://example.com", out: "jabber://example.com", err: true},
	// Normalization example:
	{in: "hTTp://subSUB.sub.EXAMPLE.COM///x//////y///foo.mp3", out: "http://subsub.sub.example.com/x/y/foo.mp3"},
	// ..more robust test cases covered by Purell
}

func TestNormalizeString(t *testing.T) {
	for i, tt := range tests {
		if got, err := urlx.NormalizeString(tt.in); got != tt.out || (err != nil && !tt.err) || (err == nil && tt.err) {
			t.Errorf(`%v. NormalizeUrl("%v") = "%v", err=%v, want "%v", err=%v`, i, tt.in, got, err, tt.out, tt.err)
		}
	}
}
