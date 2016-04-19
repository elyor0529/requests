requests
========
[![GoDoc](https://godoc.org/github.com/jochasinga/requests?status.svg)](https://godoc.org/github.com/jochasinga/requests)   [![Build Status](https://drone.io/github.com/jochasinga/requests/status.png?style=flat)](https://drone.io/github.com/jochasinga/requests/latest)   [![Coverage Status](https://coveralls.io/repos/github/jochasinga/requests/badge.svg?branch=master)](https://coveralls.io/github/jochasinga/requests?branch=master)   [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software)

Go HTTP Requests for Rodents (◕ᴥ◕)

Why Another HTTP Package?
-------------------------
Go's very own `net/http` has it all for making HTTP requests. However, `requests` wants to help
make REST calls more declarative. This means you don't have to deal with `net/http` types
and those types headaches are usually made of (such as `io.Reader` and `io.ReadCloser`).
It is safe for all [rodents](http://www.styletails.com/wp-content/uploads/2014/06/guinea-pig-booboo-lieveheersbeestje-2.jpg), not just gophers.

Install
-------

```bash

$ go get github.com/jochasinga/requests

```

Examples
--------

(*Error handling is omitted for brevity.*)

Send a common GET request just like you'd do with `http.Get`.

```go
import (
	rq "github.com/jochasinga/requests"
)

func main() {
	res, err := rq.Get("http://httpbin.org/get")
	fmt.Println(res.StatusCode)  // 200
}
```

You can send an additional data with your request.

```go

data := map[string][]string{"foo": ["bar", "baz"]}
res, err := rq.Get("http://httpbin.org/get", data)

```

To pass a basic auth with your request, just create an auth map like this:

```go

auth := map[string]string{"user": "password"}
res, err := rq.Get("http://httpbin.org/get", nil, auth)

```

You can create a custom header the same way.

```go

header := map[string][]string{"Content-Type": []{"application/json"}}
res, err := rq.Get("http://httpbin.org/get", nil, nil, header)

```

Data can be an array, slice, map, or struct.

```go

data := struct {
	foo []string
}{[]string{ "bar", "baz" }}

res, _ := rq.Post("https://httpbin.org/post", "application/json", data)

```

`GetAsync` transparently returns a channel, on which you can wait for the response.

```go

timeout := time.Duration(1) * time.Second
resChan, _ := rq.GetAsync("https://golang.org", nil, nil, timeout)

// Do some other things

res := <-resChan
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

You can also check out [relay](https://github.com/jochasinga/relay), which is
useful for testing timeouts in HTTP requests.

requests.Response
-----------------
Response returned is a special kind of `*requests.Response`, which wraps
around `*http.Response` to provide more methods such as `Text()`, `Content()`,
and `JSON()`.


### Text() ###

`Text()` returns the response's body as a string. It returns an empty string
if the body is `nil`.

```go
text := res.Text()

```

### Content() ###

`Content()` returns the response's body as `[]byte`. It returns an empty `[]byte`
if the body is `nil`.

```go

content := res.Content()

```

### JSON() ###

If one of the `Content-Type` is set to `application/json` in the response's header,
you can get the response as a JSON `[]byte`, else it returns an empty `[]byte`.

```go

jsn := res.JSON()

```

Donate
------
I am currently unemployed, working on this project alongside some others meeting
ends and expecting first child. If you find this project useful in any way or would
like to support, please [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software) me.
