// Copyright 2016 Jo Chasinga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package requests is HTTP requests made simple for rodents. It is inspired partly
by the awesomeness of HTTP request libraries in dynamic languages like Python
and Javascript. It is crafted for all rodents (and even cats), not just Gophers.

Requests is built around Go's standard http package, and in fact it is encourage
to dive in and use it to learn Go, as requests is built as a tool of convenience
than a wholely transparent one.

To send a common GET request just like you'd do with `http.Get`.

        import (
        	rq "github.com/jochasinga/requests"
        )

        func main() {
        	res, err := rq.Get("http://httpbin.org/get")
        	fmt.Println(res.StatusCode)  // 200
        }

You can send an additional data, auth and custom header with your request.


        data := map[string][]string{"foo": ["bar", "baz"]}
        auth := map[string]string{"user": "password"}
        header := map[string][]string{
                "Content-Type": []string{"application/json"},
        }

        res, err := rq.Get("http://httpbin.org/get", data, auth, header)

The data can be a map or struct.

        data1 := map[string][]string{"foo": []string{"bar", "baz"}}
        data2 := struct {
                Foo []string `json:"foo"`
        }{
                []string{"bar", "baz"},
        }

        data := map[string][]interface{}{
		"combined": {data1, data2},
	}

        res, err := rq.Post("http://httpbin.org/post", "application/json", data)

You can asynchronously wait on a GET response with `GetAsync`.

        timeout := time.Duration(1) * time.Second
        resChan, _ := rq.GetAsync("http://httpbin.org/get", nil, nil, timeout)

        // Do some other things

        res := <-resChan
        fmt.Println(res.StatusCode)  // 200

Requests is an ongoing project. Any contribution is whole-heartedly welcomed.

*/
package requests
