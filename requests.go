// Package requests provide useful and declarative methods for RESTful HTTP requests.
package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// HTTPResponse is an interface type implemented by *Response and *http.Response.
type HTTPResponse interface {
	Cookies() []*http.Cookie
	Location() (*url.URL, error)
	ProtoAtLeast(major, minor int) bool
	Write(w io.Writer) error
	Content() []byte
	Text() string
	JSON() []byte
}

// Response is a *http.Response and implements HTTPResponse.
type Response struct {
	*http.Response
}

// Text returns the response's body as string.
func (r *Response) Text() (body string) {

	if r.Body != nil {
		bodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		body = string(bodyByte)
	}

	return
}

// Content returns the response's body as []byte.
func (r *Response) Content() (content []byte) {

	if r.Body != nil {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()
		content = b
	}

	return
}

// JSON returns a JSON []byte from the response's body, provided the
// Content-Type is set as "application/json" in the response's header.
func (r *Response) JSON() (jsn []byte) {

	for _, arg := range r.Header["Content-Type"] {
		if arg == "application/json" {
			jsn = r.Content()
			return
		}
	}

	if len(r.Header["Content-Type"]) <= 0 {
		return
	}

	return
}

// Get sends a HTTP GET request to the provided URL with the data and basic authorization
// maps or structs. It returns *http.Response on success or error.
// Get accepts arguments in the followering order--urlStr, body, auth, and header.
// Only the urlStr is needed to send a request.
func Get(urlStr string, args ...interface{}) (*Response, error) {

	results, err := marshalGet(args)
	if err != nil {
		return nil, err
	}

	// Body
	body := results["body"]

	var bodyStream io.ReadCloser

	if len(body) > 0 {
		bodyStream = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	req, err := http.NewRequest("GET", urlStr, bodyStream)
	if err != nil {
		return nil, err
	}

	// auth
	auth := results["auth"]

	var authData map[string]interface{}

	if auth != nil {
		err = json.Unmarshal(auth, &authData)
		if err != nil {
			return nil, err
		}
		for usr, pw := range authData {
			password, ok := pw.(string)
			if !ok {
				return nil, err
			}

			req.SetBasicAuth(usr, password)
		}
	}

	// WATCH: header
	header := results["header"]

	var headerData http.Header

	if header != nil {
		err = json.Unmarshal(header, &headerData)
		if err != nil {
			return nil, err
		}
		req.Header = headerData
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Wrap *http.Response with *Response
	response := &Response{
		Response: resp,
	}

	return response, nil
}

// GetAsync sends a HTTP GET request to the provided URL with data and authorization
// maps or structs. It returns a chan *http.Response immediately.
//func (r *Requests) GetAsync(url string, data, auth interface{}, timeout time.Duration) (chan *http.Response, error) {
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

// Post sends a HTTP POST request to the provided URL, and encode the data according to
// the appropriate bodyType.
//func (r *Requests) Post(url, bodyType string, data interface{}) (*http.Response, error) {
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
