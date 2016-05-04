![blue gopher](http://i.imgur.com/hHNgf71.png?1)

requests
========
[![GoDoc](https://godoc.org/github.com/jochasinga/requests?status.svg)](https://godoc.org/github.com/jochasinga/requests)   [![Build Status](https://travis-ci.org/jochasinga/requests.svg?branch=master)](https://travis-ci.org/jochasinga/requests)   [![Coverage Status](https://coveralls.io/repos/github/jochasinga/requests/badge.svg?branch=master)](https://coveralls.io/github/jochasinga/requests?branch=master)   [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software)

Functional HTTP Requests in Go.

Introduction
------------
requests is a minimal, atomic, and functional way of making HTTP requests.
It is safe for all [rodents](http://www.styletails.com/wp-content/uploads/2014/06/guinea-pig-booboo-lieveheersbeestje-2.jpg), not just Gophers.

### Functional Options
requests employs functional options as parameters, this approach being
idiomatic, clean, and [makes a friendly, extensible API](http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).
This pattern is adopted after feedback from the Go community.

```go

jsontype := func(r *requests.Request) {
        r.Header.Add("content-type", "application/json")
}
res, err := requests.Get("http://example.com", jsontype)

```

### Embedded Standard Types
requests uses custom [Request](https://godoc.org/github.com/jochasinga/requests#Request)
and [Response](https://godoc.org/github.com/jochasinga/requests#Response) types
to embed standard `http.Request`, `http.Response`, and `http.Client`
in order to insert helper methods, make it easy to configure options atomically,
and [handle asynchronous errors using a special field](#handling-async-errors).

```go

timeout := func(r *requests.Request) {
        r.Timeout = time.Duration(5) * time.Second
}
res, err := requests.Get("http://example.com", timeout)
if err != nil {
        panic(err)
}

// Helper method
htmlStr := res.String()

```

See [Types and Methods](#types-and-methods) for more information.

### Asynchronous APIs
requests provides wrapper around sending an HTTP request in a goroutine , namely `requests.GetAsync` and  `requests.PostAsync`. Both return a channel on which a
`*requests.Response` can be waited on.

```go

rc, _ := requests.GetAsync("http://httpbin.org/get")
res := <-rc
// Handle connection errors.
if res.Error != nil {
        panic(res.Error)
}
// Helper method
content := res.Bytes()

```

See [Handling Async Errors](#handling-async-errors) for more information on how to handle connection errors from the goroutine.

Install
-------

```bash

go get github.com/jochasinga/requests

```

Testing
-------
requests uses Go standard `testing` package. Simple run this in the project's directory:

```bash

go test -v -cover

```

Examples
--------
### `requests.Get`
Sending a basic GET request is straightforward.

```go

res, err := requests.Get("http://httpbin.org/get")
if err != nil {
        panic(err)
}
fmt.Println(res.StatusCode)  // 200

```

To send additional data, such as a query parameter, or set basic authorization header
or content type, use functional options.

```go

// Add a query parameter.
addFoo := func(r *requests.Request) {
        r.Params.Add("foo", "bar")
}

// Set basic username and password.
setAuth := func(r *requests.Request) {
        r.SetBasicAuth("user", "pass")
}

// Set the Content-Type.
setMime := func(r *requests.Request) {
        r.Header.Add("content-type", "application/json")
}

// Pass as parameters to the function.
res, err := requests.Get("http://httpbin.org/get", addFoo, setAuth, setMime)

```

Or everything goes into one functional option.

```go

opts := func(r *requests.Request) {
        r.Params.Add("foo", "bar")
        r.SetBasicAuth("user", "pass")
        r.Header.Add("content-type", "application/json")
}
res, err := requests.Get("http://httpbin.org/get", opts)

```

### `requests.GetAsync`
After parsing all the options, spawn a goroutine to send a GET request and return `<-chan *Response` right away on which you response can be waited.

```go

timeout := func(r *requests.Request) {
        r.Timeout = time.Duration(5) * time.Second
}

rc, err := requests.GetAsync("http://golang.org", timeout)
if err != nil {
        panic(err)
}

// Do other things...

// Block and wait
res := <-rc

// Handle a "reject" with Error field.
if res.Error != nil {
	panic(res.Error)
}
fmt.Println(res.StatusCode)  // 200

```

Alternatively, `select` can be used to poll channels asynchronously.

```go

res1, _ := requests.GetAsync("http://google.com")
res2, _ := requests.GetAsync("http://facebook.com")
res3, _ := requests.GetAsync("http://docker.com")

for i := 0; i < 3; i++ {
        select {
    	case r1 := <-res1:
    		fmt.Println(r1.StatusCode)
    	case r2 := <-res2:
    		fmt.Println(r2.StatusCode)
    	case r3 := <-res3:
    		fmt.Println(r3.StatusCode)
    	}
}

```

### `requests.Post`
Send POST requests with specific `bodyType` and `body`.

```go

res, err := requests.Post("https://httpbin.org/post", "image/jpeg", &buf)

```

It also accepts variadic number of functional options:

```go

notimeout := func(r *requests.Request) {
        r.Timeout = 0
}
res, err := requests.Post("https://httpbin.org/post", "application/json", &buf, notimeout)

```

### `requests.PostAsync`
An asynchronous counterpart of `requests.Post`. Works similar to `requests.GetAsync`.

### `requests.PostJSON`
Encode your map or struct data as JSON and set `bodyType`
to `application/json` implicitly.

```go

first := map[string][]string{
        "foo": []string{"bar", "baz"},
}
second := struct {
        Foo []string `json:"foo"`
}{[]string{"bar", "baz"}}

payload := map[string][]interface{}{
        "twins": {first, second}
}

res, err := requests.PostJSON("https://httpbin.org/post", payload)

```
Types and Methods
-----------------
### `requests.Request`
It has embedded types `*http.Request` and `*http.Client`, making it an atomic
type to pass into a functional option.
It also contains field `Params`, which has the type `url.Values`. Use this field
to add query parameters to your URL. Currently, parameters in `Params` will replace
all the existing query string in the URL.

```go

addParams := func(r *requests.Request) {
        r.Params = url.Values{
	        "name": {"Ava", "Sanchez", "Poco"},
        }
}

// "q=cats" will be replaced by the new query string
res, err := requests.Get("https://httpbin.org/get?q=cats", addParams)

```

### `requests.Response`
It has embedded type `*http.Response` and provides extra byte-like helper methods
such as:
+ `Len() int`
+ `String() string`
+ `Bytes() []byte`
+ `JSON() []byte`

These methods will return an equivalent of `nil` for each return type if a
certain condition isn't met. For instance:

```go

res, _ := requests.Get("http://somecoolsite.io")
fmt.Println(res.JSON())

```

If the response from the server does not specify `Content-Type` as "application/json",
`res.JSON()` will return an empty bytes slice. It does not panic if the content type
is empty.

Another helper method, `ContentType()`, is used to get the media type in the
response's header.

```go

mime, _, err := res.ContentType()
if mime != "application/json" {
        fmt.Printf("res.JSON() returns empty %v", res.JSON())
}

```

### Handling Async Errors
`requests.Response` also has an `Error` field which will contain any error
caused in the goroutine within `requests.GetAsync` and carries it downstream
to the main goroutine for proper handling (Think `reject` in Promise but more
straightforward in Go-style).

```go
rc, _ := requests.GetAsync("http://www.docker.io")
res := <-rc
if res.Error != nil {
	panic(res.Error)
}
fmt.Println(res.StatusCode)
```

`Response.Error` is default to `nil` when there is no error or when the response
is received from a synchronous `Get`, since the error is already returned at the
function's level.

At this point `GetAsync` does not offer much more than a GET request
in a goroutine. Go is already a language built for async first-hand, and
in my opinion it does not need any extra async implementation. `GetAsync`
decidedly returns a channel, Go's native way of communication, instead
of a layer of abstraction like a Promise.

**requests** will try to be thin. You're free to pick and choose the method
that suite you best.

HTTP Test Servers
-----------------
Check out my other project [relay](https://github.com/jochasinga/relay),
useful test proxies and round-robin switchers for end-to-end HTTP tests.

Disclaimer
----------
This project is very young, but it is growing everyday since I am currently
working on this project alongside some other ideas unemployed. To support my
ends in NYC and help me push commits, please consider [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software) to fuel me with quality 🍵 or 🌟 this repo for spiritual octane.    
Reach me at [@jochasinga](http://twitter.com/jochasinga).
