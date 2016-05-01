requests
========
[![GoDoc](https://godoc.org/github.com/jochasinga/requests?status.svg)](https://godoc.org/github.com/jochasinga/requests)   [![Build Status](https://drone.io/github.com/jochasinga/requests/status.png?style=flat)](https://drone.io/github.com/jochasinga/requests/latest)   [![Coverage Status](https://coveralls.io/repos/github/jochasinga/requests/badge.svg?branch=master)](https://coveralls.io/github/jochasinga/requests?branch=master)   [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software)

Go HTTP Requests for Rodents (◕ᴥ◕)

Introduction
------------
**requests** is a minimal, atomic and expressive way of making HTTP requests.
It is inspired partly by the HTTP request libraries in other dynamic languages
like Python and Javascript. It is safe for all [rodents](http://www.styletails.com/wp-content/uploads/2014/06/guinea-pig-booboo-lieveheersbeestje-2.jpg), not just Gophers.

#### Functional Options
Passing first-class functions as optional parameters to another
function is idiomatic, clean, and [makes a friendly, extensible API](http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).
This pattern is adopted after feedbacks from the Go community.

```go

jsontype := func(r *requests.Request) {
	r.Header.Add("content-type", "application/json")
}
res, err := requests.Get("http://example.com", jsontype)

```

#### Embedded `http.Request` and `http.Response`
requests use [requests.Request](https://godoc.org/github.com/jochasinga/requests#Request)
and [requests.Response](https://godoc.org/github.com/jochasinga/requests#Response)
in order to insert helper methods and fields, make it easy to
configure options atomically, and handle asynchronous errors (See [Types and Methods](#types-and-methods)).

```go

timeout := func(r *requests.Request) {
	r.Timeout = time.Duration(5) * time.Second
}
res, err := requests.Get("http://example.com", timeout)
if err != nil {
	panic(err)
}
// helper methods
htmlStr := res.String()

```

Install
-------

```bash
go get github.com/jochasinga/requests
```

Examples
--------
#### GET requests
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

// Add a query parameter
addFoo := func(r *requests.Request) {
	r.Params.Add("foo", "bar")
}

// Set basic username and password
setAuth := func(r *requests.Request) {
	r.SetBasicAuth("user", "pass")
}

// Set the Content-Type
setMime := func(r *requests.Request) {
	r.Header.Add("content-type", "application/json")
}

// Pass as parameters to the function
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

#### Asynchronous GET
`requests.GetAsync` returns a receive-only channel `<-chan *Response`, on which
you response can be waited.

```go

timeout := func(r *requests.Request) {
	r.Timeout = time.Duration(5) * time.Second
}

rc, err := requests.GetAsync("http://golang.org", timeout)
if err != nil {
        panic(err)
}

// Do other things

// Block and wait for the response
res := <-rc
if res.Error != nil {
	panic(resp.Error)
}
fmt.Println(res.StatusCode)  // 200

```

`requests.Response` has an `Error` field that carries any error caused by
the internal goroutine to the main one so it can be handled.

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

> TODO: `requests.Pool` coming soon.

#### POST requests
`requests.Post` is used to send POST requests

```go

res, err := requests.Post("https://httpbin.org/post", "image/jpeg", &buf)
```

It also accept variadic number of functional options:

```go

notimeout := func(r *requests.Request) {
        r.Timeout = 0
}
res, err := requests.Post("https://httpbin.org/post", "application/json", &buf, notimeout)

```

`requests.PostJSON` marshals your data as JSON and set `bodyType` to
`application/json` implicitly.

```go

first := map[string][]string{"foo": []string{"bar", "baz"}}
second := struct {Foo []string `json:"foo"`}{[]string{"bar", "baz"}}
payload := map[string][]interface{}{"twins": {first, second}}

res, err := requests.PostJSON("https://httpbin.org/post", data)

```

Types and Methods
-----------------
#### `requests.Response`
Provides extra byte-like methods such as:
+ `Len() int`
+ `String() string`
+ `Bytes() []byte`
+ `JSON() []byte`

These methods will return an equivalent of `nil` for each return type if a
certain condition isn't met. For instance:

```go
resp, _ := requests.Get("http://somecoolsite.io")
fmt.Println(resp.JSON())
```

If the response from the server does not specify `Content-Type` as "application/json",
`resp.JSON()` will return an empty bytes slice. It panics if an attempt to parse
the media type fails.

Another helper method, `ContentType()`, is a shortcut for

```go

mime.ParseMediaType(request.Header.Get("content-type"))

```

#### `Handling Async Error`
`requests.Response` contains an `Error` field which carries any error
caused in the goroutine within `requests.GetAsync` downstream to the main
goroutine. It is in some way like `reject` in Promise.

```go
rc, _ := requests.GetAsync("http://www.docker.io")
res := <-rc
if res.Error != nil {
	panic(res.Error)
}
fmt.Println(res.StatusCode)
```

`Response.Error` is default to `nil` when there is no error or when the response
is received from a synchronous `Get`, since the error is already return at the
function's level.

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
