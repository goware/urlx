package urlx_test

import (
	"fmt"
	"testing"

	"github.com/goware/urlx"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		// Error out on missing host:
		{in: "", err: true},
		{in: "/", err: true},
		{in: "//", err: true},

		// Test schemes:
		{in: "http://example.com", out: "http://example.com"},
		{in: "HTTP://example.com", out: "http://example.com"},
		{in: "https://example.com", out: "https://example.com"},
		{in: "HTTPS://example.com", out: "https://example.com"},
		{in: "ssh://example.com:22", out: "ssh://example.com:22"},
		{in: "jabber://example.com:5222", out: "jabber://example.com:5222"},

		// Leading double slashes (any scheme) defaults to http:
		{in: "//example.com", out: "http://example.com"},

		// Empty scheme defaults to http:
		{in: "http://localhost", out: "http://localhost"},
		{in: "localhost", out: "http://localhost"},
		{in: "LOCALHOST", out: "http://localhost"},
		{in: "user.local", out: "http://user.local"},
		{in: "127.0.0.1", out: "http://127.0.0.1"},
		{in: "[2001:db8:a0b:12f0::1]", out: "http://[2001:db8:a0b:12f0::1]"},
		{in: "[2001:db8:a0b:12f0::80]", out: "http://[2001:db8:a0b:12f0::80]"},

		// Keep the port on non-matching scheme:
		{in: "http://localhost:80", out: "http://localhost"},
		{in: "localhost:8080", out: "http://localhost:8080"},
		{in: "[2001:db8:a0b:12f0::1]:8080", out: "http://[2001:db8:a0b:12f0::1]:8080"},
		{in: "[2001:db8:a0b:12f0::80]:80", out: "http://[2001:db8:a0b:12f0::80]"},

		// Suppress the port on matching scheme:
		{in: "localhost:80", out: "http://localhost"},
		{in: "127.0.0.1:80", out: "http://127.0.0.1"},
		{in: "[2001:db8:a0b:12f0::1]:80", out: "http://[2001:db8:a0b:12f0::1]"},

		// Test domains, subdomains etc.:
		{in: "example.com", out: "http://example.com"},
		{in: "1.example.com", out: "http://1.example.com"},
		{in: "subsub.sub.example.com", out: "http://subsub.sub.example.com"},

		// Test userinfo:
		{in: "user@example.com", out: "http://user@example.com"},
		{in: "user:passwd@example.com", out: "http://user:passwd@example.com"},
		{in: "https://user:passwd@subsub.sub.example.com", out: "https://user:passwd@subsub.sub.example.com"},

		// Lowercase scheme and host by default. Let net/url normalize URL by default:
		{in: "hTTp://subSUB.sub.EXAMPLE.COM/x//////y///foo.mp3?c=z&a=x&b=y#t=20", out: "http://subsub.sub.example.com/x//////y///foo.mp3?c=z&a=x&b=y#t=20"},
		{in: "example.com///////////////", out: "http://example.com/"},
	}

	for _, tt := range tests {
		url, err := urlx.Parse(tt.in)
		if err != nil {
			if !tt.err {
				t.Errorf(`"%s": unexpected error \"%v\"`, tt.in, err)
			}
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`"%s": expected error`, tt.in)
			continue
		}
		if url.String() != tt.out {
			t.Errorf(`"%s": got "%s", want "%v"`, tt.in, url, tt.out)
		}
	}
}

func TestURLNormalize(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		// Remove unnecessary host dots:
		{in: "http://.example.com/index.html", out: "http://example.com/index.html"},
		// Purell bugs? They claim this works..
		//{in: "http://..example..com../index.html", out: "http://example.com/index.html"},

		// Remove default port:
		{in: "http://example.com:80/index.html", out: "http://example.com/index.html"},

		// Remove duplicate slashes.
		{in: "http://example.com///x//////y///index.html", out: "http://example.com/x/y/index.html"},

		// Remove unnecesary dots from path:
		{in: "http://example.com/./x/y/z/../index.html", out: "http://example.com/x/y/index.html"},

		// Sort query:
		{in: "http://example.com/index.html?c=z&a=x&b=y", out: "http://example.com/index.html?a=x&b=y&c=z"},

		// Leave fragment as is:
		{in: "http://example.com/index.html#t=20", out: "http://example.com/index.html#t=20"},

		// README example:
		{in: "localhost:80///x///y/z/../././index.html?b=y&a=x#t=20", out: "http://localhost/x/y/index.html?a=x&b=y#t=20"},

		// ..more robust test cases covered by Purell
	}

	for _, tt := range tests {
		u, _ := urlx.Parse(tt.in)
		url, err := u.Normalize()
		if err != nil {
			if !tt.err {
				t.Errorf(`%v: unexpected error \"%v\"`, tt.in, err)
			}
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`%v: expected error`, tt.in)
			continue
		}
		if url != tt.out {
			t.Errorf(`%v: got "%v", want "%v"`, tt.in, url, tt.out)
		}
	}
}

func TestURLResolve(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		{in: "localhost", out: "127.0.0.1"},
		{in: "google.com"},
		{in: "some.weird.hostname.example.com", err: true},
	}

	for _, tt := range tests {
		u, _ := urlx.Parse(tt.in)
		ip, err := u.Resolve()
		if !tt.err && err != nil {
			t.Errorf(`%v: unexpected error \"%v\"`, tt.in, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`%v: expected error`, tt.in)
		}
		if tt.out != "" && tt.out != fmt.Sprint(ip) {
			t.Errorf(`%v: got "%v", want "%v"`, tt.in, ip, tt.out)
		}
	}

}
