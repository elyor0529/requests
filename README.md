requests
========
Go HTTP Requests for Rodents (◕ᴥ◕)

Install
-------

```bash

$ go get github.com/jochasinga/requests

```

Examples
--------

Send a request to a URL and wait for the response.

```go

auth := map[string]string{ "user" : "pass" }
res, _ := requests.Get("https://golang.org", "", auth)

fmt.Println(res.StatusCode)  // 200

```
Response returned is just a normal `*http.Response`

```go

buf := new(bytes.Buffer)
_, _ = buf.ReadFrom(res.Body)
fmt.Println(buf.String())    // Print response's Body

```

Asynchronous call will transparently return a channel, on which you can wait for the response.

```go

resChan, _ := requests.GetAsync("https://golang.org", "", nil, 2)

// Do some other things

res := <-resChan
fmt.Println(res.StatusCode)  // 200

```
Or use `select` to poll many channels asynchronously

```go

res1, _ := requests.GetAsync("http://google.com", "", nil, 1)
res2, _ := requests.GetAsync("http://facebook.com", "", nil, 1)
res3, _ := requests.GetAsync("http://docker.com", "", nil, 1)

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





