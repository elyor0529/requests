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
        	"github.com/jochasinga/requests"
        )

        func main() {
        	res, err := requests.Get("http://httpbin.org/get")
        	fmt.Println(res.StatusCode)  // 200
        }

You can send an additional data, auth and custom header with your request.


        data := map[string][]string{"foo": ["bar", "baz"]}
        auth := map[string]string{"user": "password"}
        header := map[string][]string{
                "Content-Type": []string{"application/json"},
        }

        res, err := requests.Get("http://httpbin.org/get", data, auth, header)

The data can be a map or struct (anything JSON-marshalable).

        data1 := map[string][]string{"foo": []string{"bar", "baz"}}
        data2 := struct {
                Foo []string `json:"foo"`
        }{[]string{"bar", "baz"}}

        data := map[string][]interface{}{
		"combined": {data1, data2},
	}

        res, err := requests.Post("http://httpbin.org/post", "application/json", data)

You can asynchronously wait on a GET response with `GetAsync`.

        timeout := time.Duration(1) * time.Second
        resChan, _ := requests.GetAsync("http://httpbin.org/get", nil, nil, timeout)

        // Do some other things

        res := <-resChan
        fmt.Println(res.StatusCode)  // 200

The response returned has the type *requests.Response which embeds *http.Response type
to provide more buffer-like methods such as Len(), String(), Bytes(), and JSON().

        // Len() returns the body's length.
        var len int = res.Len()

        // String() returns the body as a string.
        var text string = res.String()

        // Bytes() returns the body as bytes.
        var content []byte = res.Bytes()

        // JSON(), like Bytes() but returns an empty `[]byte` unless `Content-Type`
        // is set to `application/json` in the response's header.
        var jsn []byte = res.JSON()

These special methods use bytes.Buffer under the hood, thus unread portions of data
are returned. Make sure not to read from the response's body beforehand.

Requests is an ongoing project. Any contribution is whole-heartedly welcomed.

*/
package requests
