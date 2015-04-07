# URLx
[Golang](http://golang.org/) pkg for URL parsing and normalization. The parsing behavior is slightly different from `net/url`, see Usage.

[![GoDoc](https://godoc.org/github.com/goware/urlx?status.png)](https://godoc.org/github.com/goware/urlx)
[![Travis](https://travis-ci.org/goware/urlx.svg?branch=master)](https://travis-ci.org/goware/urlx)

## Usage

```go
import "github.com/goware/urlx"
```

### urlx.Parse()

Parse parses raw URL string into the URL struct. It uses `net/url.Parse()`
internally, but it slightly changes it's behavior:

1. It forces the default scheme and port.
2. It favors absolute paths over relative ones, thus `"example.com"` is
   parsed into `url.Host` instead of into `url.Path`.
3. It splits `Host:Port` into separate fields by default.

```go
url, _ := urlx.Parse("example.com")

// url.Scheme == "http"
// url.Port == "80"
// url.Host == "example.com"

fmt.Print(url)
// Prints http://example.com:80
```

### url.Normalize()

Normalize returns normalized URL string.
Behavior:

1. Remove unnecessary host dots.
2. Remove default port (`http://localhost:80` becomes `http://localhost`).
3. Remove duplicate slashes.
4. Remove unnecessary dots from path.
5. Sort query parameters.
6. Decode host IP into decimal numbers.
7. Handle escape values.

```go
url, _ := urlx.Parse("localhost:80///x///y/z/../././index.html?b=y&a=x#t=20")
normalized, _ := url.Normalize()

fmt.Print(normalized)
// Prints http://localhost/x/y/index.html?a=x&b=y#t=20
```

### url.Resolve()

Resolve resolves the URL host to its IP address.

```go
url, _ := urlx.Parse("localhost")
ip, _ := url.Resolve()

fmt.Print(ip)
// Prints 127.0.0.1
```

## License
URLx is licensed under the [MIT License](./LICENSE).
