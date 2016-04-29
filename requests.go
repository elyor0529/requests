// Package requests provide useful and declarative methods for
// RESTful HTTP requests.
package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// GetFunc sends a HTTP GET request to the provided url with the
// functional options to add query paramaters, headers, timeout, etc.
//
//     var addMimeType = func(r *Request) {
//             r.Header.Add("content-type", "application/json")
//     }
//
//     resp, err := requests.Get("http://httpbin.org/get", addMimeType)
//
func GetFunc(urlStr string, options ...func(*Request)) (*Response, error) {
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
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}
	// Wrap *http.Response with *Response
	response := &Response{Response: resp}
	return response, nil
}

// GetAsync sends a HTTP GET request to the provided URL with
// data and authorization maps or structs. It returns a chan
// *http.Response immediately.
func GetAsync(url string, data, auth interface{}, timeout time.Duration) (chan *http.Response, error) {
	results, err := marshalData(data, auth)
	if err != nil {
		return (chan *http.Response)(nil), err
	}
	dat, aut := results["data"], results["auth"]
	dataReadCloser := ioutil.NopCloser(bytes.NewBuffer(dat))
	req, err := http.NewRequest("GET", url, dataReadCloser)
	if err != nil {
		return (chan *http.Response)(nil), err
	}
	var authData map[string]interface{}
	if err = json.Unmarshal(aut, &authData); err != nil {
		return (chan *http.Response)(nil), err
	}
	for user, password := range authData {
		pw, ok := password.(string)
		if !ok {
			return (chan *http.Response)(nil), err
		}
		req.SetBasicAuth(user, pw)
	}
	client := &http.Client{Timeout: timeout}
	reschan := make(chan *http.Response, 1)
	go func(c chan *http.Response) error {
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		c <- res
		return nil
	}(reschan)
	return reschan, nil
}

// Post sends a HTTP POST request to the provided URL, and
// encode the data according to the appropriate bodyType.
func Post(url, bodyType string, data interface{}) (*http.Response, error) {
	dat, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataReadCloser := ioutil.NopCloser(bytes.NewBuffer(dat))
	res, err := http.DefaultClient.Post(url, bodyType, dataReadCloser)
	if err != nil {
		return (*http.Response)(nil), err
	}
	return res, nil
}
