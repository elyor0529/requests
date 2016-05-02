// Package requests provide useful and declarative methods for
// RESTful HTTP requests.
package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Head sends a HTTP HEAD request to the provided url with the
// functional options to add query paramaters, headers, timeout, etc.
//
//     addMimeType := func(r *Request) {
//             r.Header.Add("content-type", "application/json")
//     }
//
//     resp, err := requests.Head("http://httpbin.org/get", addMimeType)
//     if err != nil {
//             panic(err)
//     }
//     fmt.Println(resp.StatusCode)
//
func Head(urlStr string, options ...func(*Request)) (*Response, error) {
	req, err := http.NewRequest("HEAD", urlStr, nil)
	if err != nil {
		return nil, err
	}
	request := &Request{
		Request: req,
		Client:  &http.Client{},
		Params:  url.Values{},
	}

	// Apply options in the parameters to request.
	for _, option := range options {
		option(request)
	}
	sURL, _ := url.Parse(urlStr)
	sURL.RawQuery = request.Params.Encode()
	req.URL = sURL

	// Parse query values into r.Form
	err = req.ParseForm()
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Wrap *http.Response with *Response
	response := &Response{Response: resp}
	return response, nil
}

// Get sends a HTTP GET request to the provided url with the
// functional options to add query paramaters, headers, timeout, etc.
//
//     addMimeType := func(r *Request) {
//             r.Header.Add("content-type", "application/json")
//     }
//
//     resp, err := requests.Get("http://httpbin.org/get", addMimeType)
//     if err != nil {
//             panic(err)
//     }
//     fmt.Println(resp.StatusCode)
//
func Get(urlStr string, options ...func(*Request)) (*Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	request := &Request{
		Request: req,
		Client:  &http.Client{},
		Params:  url.Values{},
	}

	// Apply options in the parameters to request.
	for _, option := range options {
		option(request)
	}
	sURL, _ := url.Parse(urlStr)
	sURL.RawQuery = request.Params.Encode()
	req.URL = sURL

	// Parse query values into r.Form
	err = req.ParseForm()
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Wrap *http.Response with *Response
	response := &Response{Response: resp}
	return response, nil
}

// GetAsync sends a HTTP GET request to the provided URL and
// returns a <-chan *http.Response immediately.
//
//     timeout := func(r *request.Request) {
//             r.Timeout = time.Duration(10) * time.Second
//     }
//     rc, err := requests.GetAsync("http://httpbin.org/get", timeout)
//     if err != nil {
//             panic(err)
//     }
//     resp := <-rc
//     if resp.Error != nil {
//             panic(resp.Error)
//     }
//     fmt.Println(resp.String())
//
func GetAsync(urlStr string, options ...func(*Request)) (<-chan *Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	request := &Request{
		Request: req,
		Client:  &http.Client{},
		Params:  url.Values{},
	}
	for _, option := range options {
		option(request)
	}
	sURL, _ := url.Parse(urlStr)
	sURL.RawQuery = request.Params.Encode()
	req.URL = sURL
	// Parse query values into r.Form
	err = req.ParseForm()
	if err != nil {
		return nil, err
	}
	rc := make(chan *Response)
	go func() {
		resp, err := request.Client.Do(request.Request)
		// Wrap *http.Response with *Response
		response := &Response{}
		if err != nil {
			response.Error = err
			rc <- response
		}
		response.Response = resp
		rc <- response
		close(rc)
	}()
	return rc, nil
}

// Post sends a HTTP POST request to the provided URL, and
// encode the data according to the appropriate bodyType.
//
// redirect := func(r *requests.Request) {
//           r.CheckRedirect = redirectPolicyFunc
// }
// resp, err := requests.Post("https://httpbin.org/post", "image/png", &buf, redirect)
// if err != nil {
//         panic(err)
// }
// fmt.Println(resp.JSON())
//
func Post(urlStr, bodyType string, body io.Reader, options ...func(*Request)) (*Response, error) {
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	request := &Request{
		Request: req,
		Client:  &http.Client{},
		Params:  url.Values{},
	}

	// Apply options in the parameters to request.
	for _, option := range options {
		option(request)
	}
	sURL, _ := url.Parse(urlStr)
	sURL.RawQuery = request.Params.Encode()
	req.URL = sURL

	// Parse query values into r.Form and r.PostForm
	err = req.ParseForm()
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Wrap *http.Response with *Response
	response := &Response{Response: resp}
	return response, nil
}

// PostJSON aka UnsafePost! It marshals your data as JSON and set the bodyType
// to "application/json" automatically.
//
// redirect := func(r *requests.Request) {
//         r.CheckRedirect = redirectPolicyFunc
// }
//
// first := map[string][]string{"foo": []string{"bar", "baz"}}
// second := struct {Foo []string `json:"foo"`}{[]string{"bar", "baz"}}
// payload := map[string][]interface{}{"twins": {first, second}}
//
// resp, err := requests.PostJSON("https://httpbin.org/post", payload, redirect)
// if err != nil {
//         panic(err)
// }
// fmt.Println(resp.StatusCode)
//
func PostJSON(urlStr string, body interface{}, options ...func(*Request)) (*Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", urlStr, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	request := &Request{
		Request: req,
		Client:  &http.Client{},
		Params:  url.Values{},
	}
	// Apply options in the parameters to request.
	for _, option := range options {
		option(request)
	}
	sURL, _ := url.Parse(urlStr)
	sURL.RawQuery = request.Params.Encode()
	req.URL = sURL

	// Parse query values into r.Form and r.PostForm
	err = req.ParseForm()
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Wrap *http.Response with *Response
	response := &Response{Response: resp}
	return response, nil
}
