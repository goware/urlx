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
		{in: "HTTP://x.example.com", out: "http://x.example.com"},
		{in: "http://localhost", out: "http://localhost"},
		{in: "http://user.local", out: "http://user.local"},
		{in: "https://example.com", out: "https://example.com"},
		{in: "HTTPS://example.com", out: "https://example.com"},
		{in: "ssh://example.com:22", out: "ssh://example.com:22"},
		{in: "jabber://example.com:5222", out: "jabber://example.com:5222"},

		// Leading double slashes (any scheme) defaults to http:
		{in: "//example.com", out: "http://example.com"},

		// Empty scheme defaults to http:
		{in: "localhost", out: "http://localhost"},
		{in: "LOCALHOST", out: "http://localhost"},
		{in: "localhost:80", out: "http://localhost:80"},
		{in: "localhost:8080", out: "http://localhost:8080"},
		{in: "user.local", out: "http://user.local"},
		{in: "user.local:80", out: "http://user.local:80"},
		{in: "user.local:8080", out: "http://user.local:8080"},
		{in: "127.0.0.1", out: "http://127.0.0.1"},
		{in: "127.0.0.1:80", out: "http://127.0.0.1:80"},
		{in: "127.0.0.1:8080", out: "http://127.0.0.1:8080"},
		{in: "[2001:db8:a0b:12f0::1]", out: "http://[2001:db8:a0b:12f0::1]"},
		{in: "[2001:db8:a0b:12f0::80]", out: "http://[2001:db8:a0b:12f0::80]"},

		// Keep the port even on matching scheme:
		{in: "http://localhost:80", out: "http://localhost:80"},
		{in: "http://localhost:8080", out: "http://localhost:8080"},
		{in: "http://x.example.io:8080", out: "http://x.example.io:8080"},
		{in: "[2001:db8:a0b:12f0::80]:80", out: "http://[2001:db8:a0b:12f0::80]:80"},
		{in: "[2001:db8:a0b:12f0::1]:8080", out: "http://[2001:db8:a0b:12f0::1]:8080"},

		// Test domains, subdomains etc.:
		{in: "example.com", out: "http://example.com"},
		{in: "1.example.com", out: "http://1.example.com"},
		{in: "1.example.io", out: "http://1.example.io"},
		{in: "subsub.sub.example.com", out: "http://subsub.sub.example.com"},

		// Test userinfo:
		{in: "user@example.com", out: "http://user@example.com"},
		{in: "user:passwd@example.com", out: "http://user:passwd@example.com"},
		{in: "https://user:passwd@subsub.sub.example.com", out: "https://user:passwd@subsub.sub.example.com"},

		// Lowercase scheme and host by default. Let net/url normalize URL by default:
		{in: "hTTp://subSUB.sub.EXAMPLE.COM/x//////y///foo.mp3?c=z&a=x&b=y#t=20", out: "http://subsub.sub.example.com/x//////y///foo.mp3?c=z&a=x&b=y#t=20"},

		// IDNA Punycode domains.
		// TODO: net/url escapes all the fields in String() method. Should we fix it?
		{in: "http://www.žluťoučký-kůň.cz/úpěl-ďábelské-ódy", out: "http://www.%C5%BElu%C5%A5ou%C4%8Dk%C3%BD-k%C5%AF%C5%88.cz/%C3%BAp%C4%9Bl-%C4%8F%C3%A1belsk%C3%A9-%C3%B3dy"},
		{in: "http://www.xn--luouk-k-z2a6lsyxjlexh.cz/úpěl-ďábelské-ódy", out: "http://www.xn--luouk-k-z2a6lsyxjlexh.cz/%C3%BAp%C4%9Bl-%C4%8F%C3%A1belsk%C3%A9-%C3%B3dy"},
		{in: "http://żółć.pl/żółć.html", out: "http://%C5%BC%C3%B3%C5%82%C4%87.pl/%C5%BC%C3%B3%C5%82%C4%87.html"},
		{in: "http://xn--kda4b0koi.pl/żółć.html", out: "http://xn--kda4b0koi.pl/%C5%BC%C3%B3%C5%82%C4%87.html"},

		// IANA TLDs.
		// TODO: net/url escapes all the fields in String() method. Should we fix it?
		{in: "https://pressly.餐厅", out: "https://pressly.%E9%A4%90%E5%8E%85"},
		{in: "https://pressly.组织机构", out: "https://pressly.%E7%BB%84%E7%BB%87%E6%9C%BA%E6%9E%84"},

		// Some obviously wrong data:
		{in: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==", err: true},
		{in: "javascript:evilFunction()", err: true},
		{in: "otherscheme:garbage", err: true},
		{in: "<funnnytag>", err: true},
	}

	for _, tt := range tests {
		url, err := urlx.Parse(tt.in)
		if err != nil {
			if !tt.err {
				t.Errorf(`"%s": unexpected error "%v"`, tt.in, err)
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
		// Purell bug? They claim the following works..
		//{in: "http://..example..com../index.html", out: "http://example.com/index.html"},

		// Remove default port:
		{in: "http://example.com:80/index.html", out: "http://example.com/index.html"},
		{in: "localhost:80", out: "http://localhost"},
		{in: "127.0.0.1:80", out: "http://127.0.0.1"},
		{in: "[2001:db8:a0b:12f0::1]:80", out: "http://[2001:db8:a0b:12f0::1]"},

		// Empty scheme defaults to http:
		{in: "localhost", out: "http://localhost"},
		{in: "LOCALHOST", out: "http://localhost"},
		{in: "localhost:80", out: "http://localhost"},
		{in: "localhost:8080", out: "http://localhost:8080"},
		{in: "user.local", out: "http://user.local"},
		{in: "user.local:80", out: "http://user.local"},
		{in: "user.local:8080", out: "http://user.local:8080"},
		{in: "127.0.0.1", out: "http://127.0.0.1"},
		{in: "127.0.0.1:80", out: "http://127.0.0.1"},
		{in: "127.0.0.1:8080", out: "http://127.0.0.1:8080"},
		{in: "[2001:db8:a0b:12f0::1]", out: "http://[2001:db8:a0b:12f0::1]"},
		{in: "[2001:db8:a0b:12f0::1]:80", out: "http://[2001:db8:a0b:12f0::1]"},
		{in: "[2001:db8:a0b:12f0::1]:8080", out: "http://[2001:db8:a0b:12f0::1]:8080"},
		{in: "[2001:db8:a0b:12f0::80]", out: "http://[2001:db8:a0b:12f0::80]"},
		{in: "[2001:db8:a0b:12f0::80]:80", out: "http://[2001:db8:a0b:12f0::80]"},
		{in: "[2001:db8:a0b:12f0::80]:8080", out: "http://[2001:db8:a0b:12f0::80]:8080"},
		{in: "http://localhost:8080", out: "http://localhost:8080"},
		{in: "http://x.example.io:8080", out: "http://x.example.io:8080"},

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

		// Decode Punycode into UTF8.
		{in: "http://www.xn--luouk-k-z2a6lsyxjlexh.cz/úpěl-ďábelské-ódy", out: "http://www.žluťoučký-kůň.cz/%C3%BAp%C4%9Bl-%C4%8F%C3%A1belsk%C3%A9-%C3%B3dy"},
		{in: "http://xn--kda4b0koi.pl/żółć.html", out: "http://żółć.pl/%C5%BC%C3%B3%C5%82%C4%87.html"},

		// ..more robust test cases covered by Purell
	}

	for _, tt := range tests {
		u, _ := urlx.Parse(tt.in)
		url, err := urlx.Normalize(u)
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
		ip, err := urlx.Resolve(u)
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

func TestURIEncode(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		{in: "http://site.com/image/Look At Me.jpg", out: "http://site.com/image/Look%20At%20Me.jpg"},
	}

	for _, tt := range tests {
		u, err := urlx.URIEncode(tt.in)
		if !tt.err && err != nil {
			t.Errorf(`%v: unexpected error \"%v\"`, tt.in, err)
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`%v: expected error`, tt.in)
		}
		if tt.out != "" && tt.out != fmt.Sprint(u) {
			t.Errorf(`%v: got "%v", want "%v"`, tt.in, u, tt.out)
		}
	}
}
