// TODO: Break tests into multiple files.
// TODO: Change from Goconvey to standard tests.
package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/jochasinga/gtime"
	"github.com/jochasinga/relay"
	. "github.com/smartystreets/goconvey/convey"
)

func logExpectedResult(result, expected interface{}) {
	log.Printf("[RESULT]: %v\n", result)
	log.Printf("[EXPECT]: %v\n", expected)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	data := []byte(`{"foo": ["bar", "baz"]}`)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func jsonWithTypeParamHandler(w http.ResponseWriter, r *http.Request) {
	data := []byte(`{"foo": ["bar", "baz"]}`)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	html := "<html><body><h1>Blanca!</h1></body></html>"
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, html)
}

func multTypeHandler(w http.ResponseWriter, r *http.Request) {
	data := []byte(`{"foo": ["bar", "baz"]}`)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

/*************************** New GET ****************************/
var (
	fn1 = func(r *Request) {
		r.Header.Set("content-type", "application/json")
	}
	fn2 = func(r *Request) {
		r.Timeout = time.Duration(3) * time.Second
	}
	fn3 = func(r *Request) {
		r.SetBasicAuth("user", "pass")
	}
	fn4 = func(r *Request) {
		r.Params.Add("foo", "bar")
	}
	fn5 = func(r *Request) {
		r.Params.Add("name", "Ava")
	}
	getFuncArgs = [...][]func(*Request){
		{fn1},
		{fn1, fn2},
		{fn1, fn2, fn3},
		{fn1, fn2, fn3, fn4},
		{fn1, fn2, fn3, fn4, fn5},
	}
	contentTypeHandler = func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, r.Header.Get("Content-Type"))
	}
	fooHandler = func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, r.FormValue("foo"))
	}
	nameHandler = func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, r.FormValue("name"))
	}
	basicAuthHandler = func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok {
			log.Panicln("Error getting Basic Auth.")
		}
		io.WriteString(w, user+" : "+password)
	}
	getFuncTestTable = []struct {
		fn       func(*Request)
		handler  func(http.ResponseWriter, *http.Request)
		expected string
	}{
		{fn1, contentTypeHandler, "application/json"},
		{fn3, basicAuthHandler, "user : pass"},
		{fn4, fooHandler, "bar"},
		{fn5, nameHandler, "Ava"},
	}
	getFuncSyncTestTable = []struct {
		delay    int
		expected int
	}{
		{1, 2},
		{2, 4},
		{3, 6},
		{4, 8},
	}
	getFuncTimeoutTestTable = []struct {
		delay    int
		timeout  float64
		expected float64
	}{
		{1, 0.5, 0.5},
		{2, 0.5, 0.5},
		{2, 1.0, 1.0},
		{3, 1.0, 1.0},
	}
	jsonFuncTestTable = []struct {
		fn       func(http.ResponseWriter, *http.Request)
		expected []byte
	}{
		{jsonHandler, []byte(`{"foo": ["bar", "baz"]}`)},
		{jsonWithTypeParamHandler, []byte(`{"foo": ["bar", "baz"]}`)},
		{multTypeHandler, []byte(`{"foo": ["bar", "baz"]}`)},
		{htmlHandler, []byte{}},
	}
)

// Test that the returned type is always *Response.
func TestGetFuncResponseType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, fns := range getFuncArgs {
		resp, err := GetFunc(ts.URL, fns...)
		if err != nil {
			t.Error(err)
		}
		if reflect.TypeOf(resp) != reflect.TypeOf(&Response{}) {
			t.Error(err)
		}
	}
}

// Test that the request has the appropriate options.
func TestGetFuncResponseOptions(t *testing.T) {
	for _, tt := range getFuncTestTable {
		ts := httptest.NewServer(http.HandlerFunc(tt.handler))
		defer ts.Close()
		resp, err := GetFunc(ts.URL, tt.fn)
		if err != nil {
			t.Error(err)
		}
		if resp.String() != tt.expected {
			t.Error(err)
		}
	}
}

// Get should wait for the response and return
func TestGetFuncResponseTime(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, tt := range getFuncSyncTestTable {
		delay := time.Duration(tt.delay) * time.Second
		expected := time.Duration(tt.expected) * time.Second
		p := relay.NewProxy(delay, ts)
		start := time.Now()
		_, _ = GetFunc(p.URL)
		elapsed := time.Since(start)
		if elapsed <= expected {
			logExpectedResult(elapsed, expected)
			t.Error("Client returned before it should.")
		}
	}
}

// Get should wait fo the response until timed out.
func TestGetFuncResponseOnTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, tt := range getFuncTimeoutTestTable {
		delay := time.Duration(tt.delay) * time.Second
		p := relay.NewProxy(delay, ts)
		start := time.Now()
		_, err := GetFunc(p.URL, func(r *Request) {
			r.Timeout = gtime.Ftos(tt.timeout)
		})
		if err == nil {
			t.Error(errors.New("Client did not time out."))
		}
		elapsed := time.Since(start).Seconds()
		deviation := gtime.FloatTime.Seconds()
		if !(elapsed >= tt.expected-deviation || elapsed <= tt.expected+deviation) {
			logExpectedResult(elapsed, tt.expected)
			t.Error("Client returned before it should.")
		}
	}
}

func TestGetFuncResponseAsBytes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	resp, err := GetFunc(ts.URL)
	if err != nil {
		t.Error(err)
	}
	result := resp.Bytes()
	expected := []byte("Hello world!")
	if bytes.Compare(result, expected) != 0 {
		t.Error("Unexpected result.")
	}
}

func TestGetFuncResponseAsJSON(t *testing.T) {
	for _, tt := range jsonFuncTestTable {
		ts := httptest.NewServer(http.HandlerFunc(tt.fn))
		defer ts.Close()
		resp, err := GetFunc(ts.URL)
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(resp.JSON(), tt.expected) != 0 {
			t.Error(err)
		}
	}
}

func TestGetFuncResponseAsString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	resp, err := GetFunc(ts.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.String() != "Hello world!" {
		t.Error(err)
	}
}

/*************************** New GET ****************************/

var badURLs = []string{
	"://maggot.#&",
	"crap://bs.com",
	"htp://f#as3",
}

func TestGetWithBadURLs(t *testing.T) {
	for _, url := range badURLs {
		_, err := GetFunc(url)
		if err == nil {
			t.Error(err)
		}
	}
}

var getAsyncTestTable = []struct {
	delay    int
	expected int
}{
	{1, 0},
	{2, 0},
	{3, 0},
	{4, 0},
}

// GetAsync should return immediately.
func TestGetAsyncResponseTimes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	deviation := time.Duration(10) * time.Millisecond
	for _, tt := range getFuncSyncTestTable {
		delay := time.Duration(tt.delay) * time.Second
		expected := time.Duration(tt.expected)*time.Second + deviation
		p := relay.NewProxy(delay, ts)
		start := time.Now()
		_, _ = GetAsync(p.URL, nil, nil, 0)
		elapsed := time.Since(start)
		if elapsed >= expected {
			logExpectedResult(elapsed, expected)
			t.Error("Client takes too long to be asynchronous.")
		}
	}
}

// TODO: Change to standard tests
func TestGetAsyncResponseTypeAndContent(t *testing.T) {
	Convey("GIVEN the Server Handler with delay proxy", t, func() {

		var d map[string]interface{}
		timeout := time.Duration(5) * time.Second

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			msg := "Hello, requests"
			if r.Body != nil {
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
				}
				if err = json.Unmarshal(data, &d); err != nil {
					t.Error(err)
				}
			}
			if len(d) > 0 {
				for k := range d {
					msg += ", " + k
				}
			}
			if len(r.Header["Authorization"]) > 0 {
				msg += ", auth"
			}
			fmt.Fprintf(w, msg)
		}))

		latency := time.Duration(1) * time.Second
		proxy := relay.NewProxy(latency, ts)

		Convey("WITH data and auth maps", func() {
			auth := map[string]string{"user": "password"}
			data := map[string][]string{"foo": []string{"bar", "baz"}}

			rc, err := GetAsync(proxy.URL, data, auth, timeout)
			if err != nil {
				t.Error(err)
			}

			Convey("EXPECT GetAsync() to return chan *httpResponse", func() {
				So(rc, ShouldHaveSameTypeAs, make(chan *http.Response))
			})
			Convey("EXPECT GetAsync() to return correct content", func() {
				resp := <-rc
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				defer resp.Body.Close()

				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo, auth")
			})
		})

		Convey("WITH data and auth structs", func() {
			auth := struct {
				User     string
				Password string
			}{"user", "password"}

			data := struct {
				Foo []string `json:"foo"`
			}{[]string{"bar", "baz"}}

			rc, err := GetAsync(proxy.URL, data, auth, timeout)
			if err != nil {
				t.Error(err)
			}

			Convey("EXPECT GetAsync() to return chan *httpResponse", func() {
				So(rc, ShouldHaveSameTypeAs, make(chan *http.Response))
			})

			Convey("EXPECT GetAsync() to return correct content", func() {
				resp := <-rc
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				defer resp.Body.Close()

				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo, auth")
			})
		})

		Convey("WITH data and auth as nil", func() {
			rc, err := GetAsync(proxy.URL, nil, nil, timeout)
			if err != nil {
				t.Error(err)
			}

			Convey("EXPECT GetAsync() to return type chan *httpResponse", func() {
				So(rc, ShouldHaveSameTypeAs, make(chan *http.Response))
			})
			Convey("EXPECT GetAsync() to return correct content", func() {
				resp := <-rc
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				defer resp.Body.Close()

				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests")
			})
		})

		// Edge cases for errors
		Convey("WITH data as a complex number", func() {
			badData := 12i

			resp, err := GetAsync(ts.URL, badData, nil, timeout)
			Convey("EXPECT Get() to return an error", func() {
				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("WITH auth as a channel", func() {
			badAuth := make(chan complex64)

			resp, err := GetAsync(ts.URL, nil, badAuth, timeout)
			Convey("EXPECT Get() to return an error", func() {
				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("WITH malformed URL", func() {
			badURL := "://maggot.#&"

			resp, err := GetAsync(badURL, nil, nil, timeout)
			Convey("EXPECT Get() to return an error", func() {
				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Reset(func() {
			ts.Close()
			proxy.Close()
		})
	})
}

func TestPostResponseTypeAndContent(t *testing.T) {
	Convey("GIVEN the Server Handler", t, func() {
		var d map[string]interface{}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			msg := "Hello, requests"
			if r.Body != nil {
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
				}

				if err = json.Unmarshal(data, &d); err != nil {
					t.Error(err)
				}
			}

			if len(d) > 0 {
				for k := range d {
					msg += ", " + k
				}
			}

			if len(r.Header["Authorization"]) > 0 {
				msg += ", auth"
			}

			fmt.Fprintf(w, msg)
		}))

		Convey("WITH data map", func() {
			data := map[string][]string{"foo": []string{"bar", "baz"}}

			resp, err := Post(ts.URL, "application/json", data)
			if err != nil {
				t.Error(err)
			}

			Convey("EXPECT Get() to return *httpResponse", func() {
				So(resp, ShouldHaveSameTypeAs, &http.Response{})
			})
			Convey("EXPECT Get() to return correct content", func() {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				defer resp.Body.Close()

				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo")
			})

		})

		Convey("WITH data struct", func() {
			data := struct {
				Foo []string `json:"foo"`
			}{[]string{"bar", "baz"}}

			resp, err := Post(ts.URL, "application/json", data)
			if err != nil {
				t.Error(err)
			}

			Convey("EXPECT Post() to return *httpResponse", func() {
				So(resp, ShouldHaveSameTypeAs, &http.Response{})
			})

			Convey("EXPECT Post() to return correct content", func() {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				defer resp.Body.Close()

				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo")
			})
		})

		Convey("WITH data as nil", func() {
			resp, err := Post(ts.URL, "application/json", nil)
			if err != nil {
				t.Error(err)
			}

			Convey("EXPECT Post() to return type *httpResponse", func() {
				So(resp, ShouldHaveSameTypeAs, &http.Response{})
			})

			Convey("EXPECT Post() to return correct content", func() {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				defer resp.Body.Close()

				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests")
			})
		})

		// Edge cases for errors
		Convey("WITH data as a complex number", func() {
			badData := 12i
			resp, err := Post(ts.URL, "application/json", badData)

			Convey("EXPECT Post() to return an error", func() {
				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("WITH malformed URL", func() {
			badURL := "://maggot.#&"
			resp, err := Post(badURL, "application/json", nil)

			Convey("EXPECT Post() to return an error", func() {
				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Reset(func() {
			ts.Close()
		})
	})
}
