requests
========
[![GoDoc](https://godoc.org/github.com/jochasinga/requests?status.svg)](https://godoc.org/github.com/jochasinga/requests)   [![Build Status](https://drone.io/github.com/jochasinga/requests/status.png?style=flat)](https://drone.io/github.com/jochasinga/requests/latest)   [![Coverage Status](https://coveralls.io/repos/github/jochasinga/requests/badge.svg?branch=master)](https://coveralls.io/github/jochasinga/requests?branch=master)   [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software)

Go HTTP Requests for Rodents (◕ᴥ◕)

Why Another HTTP Package?
-------------------------
Go's very own `net/http` has it all for making HTTP requests. However, `requests` wants to help
make REST calls more declarative. It is safe for all [rodents](http://www.styletails.com/wp-content/uploads/2014/06/guinea-pig-booboo-lieveheersbeestje-2.jpg), not just gophers.

Install
-------

```bash

go get github.com/jochasinga/requests

```

Examples
--------

(*Error handling omitted for brevity.*)

A basic GET request:

```go
import (
	"github.com/jochasinga/requests"
)

func main() {

	res, err := requests.Get("http://httpbin.org/get")

	fmt.Println(res.StatusCode)  // 200
}
```

Additional data with request:

```go

data := map[string][]string{"foo": ["bar", "baz"]}
res, err := requests.Get("http://httpbin.org/get", data)

```

It isn't common to pass data with GET. However, in some REST endpoints,
i.e. Elasticsearch with which this package was initially written to interact,
JSON data is expected with the request as a DSL for querying.

Basic auth can also be sent as an argument:

```go

auth := map[string]string{"user": "password"}
res, err := requests.Get("http://httpbin.org/get", nil, auth)

```

For custom headers, it's best to create a `map[string][]string`
which adheres to the structure of `http.Header`:

```go

headers := map[string][]string{
	"Content-Type": {"application/json"},
	"Accept": {"text/html"},
}
res, err := requests.Get("http://httpbin.org/get", nil, nil, headers)

```

Data can be anything JSON-marshalable (map[string]interface{} or struct).

```go

data1 := map[string][]string{"foo": []string{"bar", "baz"}}
data2 := struct {
	Foo []string `json:"foo"`
}{[]string{"bar", "baz"}}
combined := map[string][]interface{}{
	"twins": {data1, data2},
}

res, err := requests.Post("https://httpbin.org/post", "application/json", combined)

```

GetAsync transparently returns a channel, on which you can wait for the response.

```go

timeout := time.Duration(1) * time.Second
rc, err := rq.GetAsync("https://golang.org", nil, nil, timeout)

// Do some other things

res := <-rc
fmt.Println(res.StatusCode)  // 200

```
Or use `select` to poll many channels asynchronously.

```go

res1, _ := rq.GetAsync("http://google.com", nil, nil, timeout)
res2, _ := rq.GetAsync("http://facebook.com", nil, nil, timeout)
res3, _ := rq.GetAsync("http://docker.com", nil, nil, timeout)

for i := 0; i < 3; i++ {
        select {
	    case r1 := <-res1:
		        fmt.Println(r1.Status)
	    case r2 := <-res2:
		        fmt.Println(r2.StatusCode)
	    case r3 := <-res3:
		        fmt.Println(r3.Body)
		}
}

```

*TODO: `requests.Pool` coming soon*

Also check out my other project [relay](https://github.com/jochasinga/relay),
which is useful for end-to-end HTTP tests.

requests.Response
-----------------
Response returned is a special kind of `*requests.Response`, which wraps
around `*http.Response` to provide more methods such as `Len(), ``String()`,
`Content()`, and `JSON()`.

### Len() ###
Returns the body's length.

```go

len := res.Len()

```

### String() ###
Returns the body as a string.

```go

text := res.String()

```

### Content() ###
Returns the body as bytes.

```go

content := res.Content()

```

### JSON() ###
Like `Content()`, but returns the body as bytes only if `Content-Type` is set
to `application/json` in the response's header. Otherwise, it returns empty `[]byte`.


```go

jsn := res.JSON()

```

These special methods use `bytes.Buffer` under the hood, so they all returns
"unread" portion of data. Make sure not to read from the data before using.

Donate
------
I am currently working on this project alongside some other ideas to meet
ends in NYC unemployed. Please consider [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software) to fuel me with proper coffee.
