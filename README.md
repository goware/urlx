# URLx
[Golang](http://golang.org/) pkg for parsing, normalization and validation of URLs. Unlike `net/url`, it favors absolute paths over relative paths (`"example.com"` is parsed as `url.Domain`, not as `url.Path`).

[![GoDoc](https://godoc.org/github.com/goware/urlx?status.png)](https://godoc.org/github.com/goware/urlx)
[![Travis](https://travis-ci.org/goware/urlx.svg?branch=master)](https://travis-ci.org/goware/urlx)

## Usage

```go
url := urlx.NormalizeString("example.com")

fmt.Println(url)
// Prints http://example.com
```

```go
url := urlx.NormalizeString("localhost:8080/ping")

fmt.Println(url)
// Prints http://localhost:8080/ping
```

## License
URLx is licensed under the [MIT License](./LICENSE).
