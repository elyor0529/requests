requests
========
[![GoDoc](https://godoc.org/github.com/jochasinga/requests?status.svg)](https://godoc.org/github.com/jochasinga/requests)   [![Build Status](https://drone.io/github.com/jochasinga/requests/status.png?style=flat)](https://drone.io/github.com/jochasinga/requests/latest)   [![Coverage Status](https://coveralls.io/repos/github/jochasinga/requests/badge.svg?branch=master)](https://coveralls.io/github/jochasinga/requests?branch=master)   [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software)

Go HTTP Requests for Rodents (◕ᴥ◕)

Why Another HTTP Package?
-------------------------
This is Go's [net/http](https://golang.org/pkg/net/http/) on steroid for making high-level HTTP requests to REST services more declarative. It is safe for all [rodents](http://www.styletails.com/wp-content/uploads/2014/06/guinea-pig-booboo-lieveheersbeestje-2.jpg) , not just gophers.

Install
-------

```bash

go get github.com/jochasinga/requests

```

Examples
--------

(*Error handling omitted for brevity.*)

### Sending GET requests

```go
import (
	"github.com/jochasinga/requests"
)

func main() {

	res, err := requests.Get("http://httpbin.org/get")

	fmt.Println(res.StatusCode)  // 200
}
```

Additional data:

```go

data := map[string][]string{"foo": ["bar", "baz"]}
res, err := requests.Get("http://httpbin.org/get", data)

```


>> Note that [it is not common to pass data with GET request](http://stackoverflow.com/questions/978061/http-get-with-request-body). However,
[some services](http://stackoverflow.com/questions/14339696/elasticsearch-post-with-json-search-body-vs-get-with-json-in-url) expect a JSON body with GET request as a [DSL](https://en.wikipedia.org/wiki/Domain-specific_language) for i.e. results querying, which is often more
convenient than passing long query parameters in the URL. It is generally a better idea
to use POST when trying to pass data to the server, and in fact I will deprecate
the body argument in the `Get()` method soon.    

Basic auth can also be sent:

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

### Sending POST requests

Data can be anything JSON-marshalable (map and struct).

```go

data1 := map[string][]string{"foo": []string{"bar", "baz"}}
data2 := struct {
	Foo []string `json:"foo"`
}{[]string{"bar", "baz"}}
data := map[string][]interface{}{
	"combined": {data1, data2},
}

res, err := requests.Post("https://httpbin.org/post", "application/json", data)

```

### Asynchronous Requests

`GetAsync` transparently returns a channel, on which you can wait for the response.

```go

timeout := time.Duration(1) * time.Second

rc, err := requests.GetAsync("https://golang.org", nil, nil, timeout)

// Do other things

res := <-rc
fmt.Println(res.StatusCode)  // 200

```
Or use `select` to poll channels asynchronously.

```go

res1, _ := requests.GetAsync("http://google.com", nil, nil, timeout)
res2, _ := requests.GetAsync("http://facebook.com", nil, nil, timeout)
res3, _ := requests.GetAsync("http://docker.com", nil, nil, timeout)

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

>> TODO: `requests.Pool` coming soon

Awesome HTTP Tests
------------------
Check out my other project [relay](https://github.com/jochasinga/relay),
which is useful for end-to-end HTTP tests.

[requests.Response](https://godoc.org/github.com/jochasinga/requests#Response)
------------------------------------------------------------------------------
The response returned has the type [*requests.Response](https://godoc.org/github.com/jochasinga/requests#Response),
which embeds [*http.Response](https://golang.org/pkg/net/http/#Response) type
to provide more buffer-like methods such as:

+ `Len()`
+ `String()`
+ `Bytes()`
+ `JSON()`

### Len()
Returns the body's length.

```go

var len int = res.Len()

```

### String() ###
Returns the body as a string.

```go

var text string = res.String()

```

### Bytes()
Returns the body as bytes.

```go

var content []byte = res.Bytes()

```

### JSON()
Like `Bytes()` but returns an empty `[]byte` unless `Content-Type` is set
to `application/json` in the response's header.

```go

var jsn []byte = res.JSON()

```

These special methods use [bytes.Buffer](https://golang.org/pkg/bytes/#Buffer)
under the hood, thus unread portions of data are returned. Make sure not to read
from the response's body beforehand.

Donate
------
I am currently working on this project alongside some other ideas to meet
ends in NYC **unemployed**. Please consider [![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=jochasinga&url=https://github.com/jochasinga/requests&title=Relay&language=English&tags=github&category=software) to fuel me with proper coffee.
