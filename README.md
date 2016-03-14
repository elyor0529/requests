requests
========
Go HTTP Requests for Rodents (◕ᴥ◕)

Example
-------

Send a request to a URL and wait for the response.

```go

res, _ := requests.Get("https://golang.org", `{"foo": "bar"}`, map[string]string{"user": "pass"})

fmt.Println(res.StatusCode)              // 200
fmt.Println(res.Header["Content-Type"])  // "application/json"
fmt.Println(res.Body)                    // {"foo": "bar"}

```

Asynchronous call will transparently return a channel

```go

resChan, _ := requests.GetAsync("https://golang.org", nil, nil)

// Do some other things

res := <-resChan
fmt.Println(res.StatusCode)  // 200

```





