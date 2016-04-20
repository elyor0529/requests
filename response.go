package requests

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// HTTPResponse is an interface type implemented by *Response and *http.Response.
type HTTPResponse interface {
	Cookies() []*http.Cookie
	Location() (*url.URL, error)
	ProtoAtLeast(major, minor int) bool
	Write(w io.Writer) error
	Bytes() []byte
	String() string
	JSON() []byte
	Len() int
}

// Response is a *http.Response and implements HTTPResponse.
type Response struct {
	*http.Response
}

// Len returns the response's body's unread portion's length,
// which is the full length provided it has not been read.
func (r *Response) Len() (len int) {

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	len = buf.Len()

	return
}

// String returns the response's body as string. Any errors
// reading from the body is ignored for convenience.
func (r *Response) String() (body string) {

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	body = buf.String()

	return
}

// Bytes returns the response's Body as []byte. Any errors
// reading from the body is ignored for convenience.
func (r *Response) Bytes() (content []byte) {

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	content = buf.Bytes()

	return
}

// JSON, like Bytes() but returns an empty []bytes "Content-Type" is set to
// "application/json" in the response's header.
func (r *Response) JSON() (jsn []byte) {

	for _, arg := range r.Header["Content-Type"] {
		if arg == "application/json" {
			jsn = r.Bytes()
			return
		}
	}

	if len(r.Header["Content-Type"]) <= 0 {
		return
	}

	return
}
