requests
========
[![GoDoc](https://godoc.org/github.com/jochasinga/requests?status.svg)](https://godoc.org/github.com/jochasinga/requests)  [![Build Status](https://drone.io/github.com/jochasinga/requests/status.png?style=flat)](https://drone.io/github.com/jochasinga/requests/latest)  [![Coverage Status](https://coveralls.io/repos/github/jochasinga/requests/badge.svg?style=flat)](https://coveralls.io/github/jochasinga/requests?branch=master)

Go HTTP Requests for Rodents (◕ᴥ◕)

Why Another HTTP Package?
-------------------------
Go's very own `net/http` has it all for making HTTP requests. However, `requests` wants to help
make REST calls more yummy and declarative. This Means you don't have to deal with `net/http` types 
and those types headaches are usually made of (such as `io.Reader` and `io.ReadCloser`).
It is safe for all [rodents](http://www.styletails.com/wp-content/uploads/2014/06/guinea-pig-booboo-lieveheersbeestje-2.jpg), not just gophers. 

Install
-------

```bash

$ go get github.com/jochasinga/requests

```

Examples
--------

Send a GET request to a URL and wait for the response.

```go

import (
	"github.com/jochasinga/requests"
)


func main() {
	rq := requests.New()
	auth := map[string]string{ "user" : "pass" }
	res, _ := rq.Get("https://golang.org", nil, auth)

	fmt.Println(res.StatusCode)  // 200
}

```

Data can be an array, slice, map, or struct.

```go

data := struct {
	foo []string
}{ []string{ "bar", "baz" } }
res, _ := rq.Post("https://httpbin.org/post", "application/json", data)

```

Response returned is just a normal `*http.Response`

```go

buf := new(bytes.Buffer)
_, _ = buf.ReadFrom(res.Body)
fmt.Println(buf.String())    // Print response's Body

```

`GetAsync` transparently returns a channel, on which you can wait for the response.

```go

timeout := time.Duration(1) * time.Second
resChan, _ := rq.GetAsync("https://golang.org", nil, nil, timeout)

// Do some other things

res := <-resChan
fmt.Println(res.StatusCode)  // 200

```
Or use `select` to poll many channels asynchronously

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

*TODO: requests.Pool coming soon*




