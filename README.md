requests
========
[![GoDoc](https://godoc.org/github.com/jochasinga/requests?status.svg)](https://godoc.org/github.com/jochasinga/requests)   [![Build Status](https://drone.io/github.com/jochasinga/requests/status.png?style=flat)](https://drone.io/github.com/jochasinga/requests/latest)   [![Coverage Status](https://coveralls.io/repos/github/jochasinga/requests/badge.svg?branch=master)](https://coveralls.io/github/jochasinga/requests?branch=master)   [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software)

Go HTTP Requests for Rodents (◕ᴥ◕)

Introduction
------------
**requests** is a minimal, atomic and expressive way of making HTTP requests.
It is safe for all [rodents](http://www.styletails.com/wp-content/uploads/2014/06/guinea-pig-booboo-lieveheersbeestje-2.jpg), not just Gophers.

Differences
-----------
#### Functional Options
Go has first-class functions, and passing them as option parameters to another
function feels idiomatic, clean, and [makes a friendly, extensible API](http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).
I adopted this pattern after many feedbacks from the Go community.

```go
res, err := requests.Get("http://example.com", func(r *requests.Request) {
	r.Header.Add("content-type", "application/json")
})
```

However, just like in Javascript, the recommended way to pass in functional
parameters is as declared variables rather than an anonymous ones.

#### Embedded `http.Request` and `http.Response`
requests use [requests.Request](https://godoc.org/github.com/jochasinga/requests#Request)
and [requests.Response](https://godoc.org/github.com/jochasinga/requests#Response)
in order to insert special helpful methods and fields and make it possible to
configure request's and client's options atomically.

```go
timeout := func(r *requests.Request) {
	r.Timeout = time.Duration(5) * time.Second
}
res, err := requests.Get("http://example.com", timeout)
if err != nil {
	panic(err)
}
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

Or everything goes into one functional option (less tedious).

```go
opts := func(r *requests.Request) {
	r.Params.Add("foo", "bar")
	r.SetBasicAuth("user", "pass")
	r.Header.Add("content-type", "application/json")
}
res, err := requests.Get("http://httpbin.org/get", opts)
```

#### Asynchronous GET
`requests.GetAsync` returns a channel `<-chan *Response`, on which you can
wait for the response.

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
if resp.Error != nil {
	panic(resp.Error)
}
fmt.Println(res.StatusCode)  // 200
```

`requests.Response` has an `Error` field that carries any error caused by
the internal goroutine to the main one so it can be handled.

Use `select` to poll channels asynchronously.

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
`requests.Post` is analogous to `http.Post`

```go
res, err := requests.Post("https://httpbin.org/post", "image/jpeg", &buf)
```

except it also accept variadic number of functional options:

```go
addProtobuf := func(r *requests.Request) {
        r.Header.Add("content-type", "application/x-protobuf")
}
res, err := requests.Post("https://httpbin.org/post", "application/json", &buf, addProtobuf)
```

`requests.PostJSON` aka "risky POST". It marshals your data as JSON and set the
bodyType to "application/json" automatically.

```go
first := map[string][]string{"foo": []string{"bar", "baz"}}
second := struct {Foo []string `json:"foo"`}{[]string{"bar", "baz"}}
payload := map[string][]interface{}{"twins": {first, second}}

timeout := func(r *requests.Request) {
	r.Timeout = time.Duration(15) * time.Second
}

res, err := requests.PostJSON("https://httpbin.org/post", data, timeout)
```

#### More on `requests.Response`
`requests.Response` provides extra byte-like methods such as:
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

If the response from the server does not specify `Content-Type` as `application/json`,
`res.JSON()` will return an empty bytes slice. It panics if an attempt to parse
the media type fails.

Another helper method, `ContentType() (string, map[string]string, error)`, is
an alias for

```go
mime.ParseMediaType(request.Header.Get("content-type"))
```

#### Handling Error from Asynchronous Calls
`requests.Response` also contains an `Error` field which carries any error
caused after a goroutine within `requests.GetAsync` is spawned over to the main
routine (Think reject in Promise).

This is how you can handle such error:

```go
rc, _ := requests.GetAsync("http://www.docker.io")
res := <-rc
if res.Error != nil {
	panic(res.Error)
}
fmt.Println(res.StatusCode)
```

Awesome HTTP Tests
------------------
Check out my other project [relay](https://github.com/jochasinga/relay),
useful servers for end-to-end HTTP tests.

Disclaimer
----------
This project is very young and is not yet production-ready, but it is growing everyday since I am currently working on this project alongside some other ideas unemployed. To support my ends in NYC
and help me push commits, please consider [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software) to fuel me with quality 🍵 or 🌟 this repo for spiritual octane.
