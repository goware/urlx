# URLx
[Golang](http://golang.org/) pkg for URL parsing and normalization.

[![GoDoc](https://godoc.org/github.com/goware/urlx?status.png)](https://godoc.org/github.com/goware/urlx)
[![Travis](https://travis-ci.org/goware/urlx.svg?branch=master)](https://travis-ci.org/goware/urlx)

## Parsing URL

```go
Parse(rawURL string) (*url.URL, error)
```

The [urlx.Parse()](https://godoc.org/github.com/goware/urlx#Parse) is compatible with the same function from [net/url](https://golang.org/pkg/net/url/#Parse) pkg, but has slightly different behavior. It enforces default scheme and favors absolute URLs over relative paths. See [GoDoc](https://godoc.org/github.com/goware/urlx#Parse) for more details.

### Example

```go
import "github.com/goware/urlx"

function main() {
    url, _ := urlx.Parse("example.com")
    // url.Scheme == "http"
    // url.Host == "example.com"

    fmt.Print(url)
    // Prints http://example.com
}
```

## Normalizing URL

```go
func Normalize(url *URL) (string, error)
```
```go
func NormalizeString(rawURL string) (string, error)
```

The [urlx.Normalize()](https://godoc.org/github.com/goware/urlx#Normalize) function normalizes the URL using the predefined subset of [Purell](https://github.com/PuerkitoBio/purell) flags. See [GoDoc](https://godoc.org/github.com/goware/urlx#Normalize) for more details.

### Example

```go
import "github.com/goware/urlx"

function main() {
    url, _ := urlx.Parse("localhost:80///x///y/z/../././index.html?b=y&a=x#t=20")
    normalized, _ := urlx.Normalize(url)

    fmt.Print(normalized)
    // Prints http://localhost/x/y/index.html?a=x&b=y#t=20
}
```

## Splitting host:port from URL

```go
func SplitHostPort(*url.URL) (host, port string, err error) 
```

The [urlx.SplitHostPort()](https://godoc.org/github.com/goware/urlx#SplitHostPort) is compatible with the same function from [net](https://golang.org/pkg/net/) pkg, but has slightly different behavior. It doesn't remove brackets from `[IPv6]` host. See [GoDoc](https://godoc.org/github.com/goware/urlx#SplitHostPort) for more details.

### Example

```go
import "github.com/goware/urlx"

function main() {
    url, _ := urlx.Parse("localhost:80")
    host, port, _ := urlx.SplitHostPort(url)

    fmt.Print(host)
    // Prints localhost

    fmt.Print(port)
    // Prints 80
}
```

## Resolving IP address from URL

```go
func Resolve(url *URL) (*net.IPAddr, error)
```
```go
func ResolveString(rawURL string) (*net.IPAddr, error)
```

The [urlx.Resolve()](https://godoc.org/github.com/goware/urlx#Resolve) is compatible with [ResolveIPAddr()](https://golang.org/pkg/net/#ResolveIPAddr) from [net](https://golang.org/pkg/net/). See [GoDoc](https://godoc.org/github.com/goware/urlx#Resolve) for more details.

### Example

```go
url, _ := urlx.Parse("localhost")
ip, _ := urlx.Resolve(url)

fmt.Print(ip)
// Prints 127.0.0.1
```

## License
URLx is licensed under the [MIT License](./LICENSE).
