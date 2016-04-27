// TODO: Break tests into multiple files.
// TODO: Change from Goconvey to standard tests.
package requests

import (
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

	"github.com/jochasinga/relay"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	bodyMap   = map[string][]string{"foo": []string{"bar", "baz"}}
	authMap   = map[string]string{"user": "password"}
	headerMap = map[string][]string{
		"Accept-Encoding": {"gzip, deflate"},
		"Accept-Language": {"en-us"},
		"Content-Type":    {"application/json"},
		"Foo":             {"Bar", "two"},
	}
	bodyHybridMap = map[string][]interface{}{
		"duplica": {bodyMap, bodyStruct},
	}
	respJSON   = []byte(`{"foo": [{"bar", "baz"}]`)
	bodyStruct = struct {
		Foo []string `json:"foo"`
	}{[]string{"bar", "baz"}}
	authStruct = struct {
		User string `json:"user"`
	}{"password"}
	headerStruct = struct {
		ContentType    []string
		AcceptEncoding []string
		AcceptLanguage []string
		Foo            []string
	}{
		[]string{"application/json"},
		[]string{"gzip, deflate"},
		[]string{"en-us"},
		[]string{"Bar", "two"},
	}
)

func logExpectedResult(result, expected interface{}) {
	log.Printf("[RESULT]: %v\n", result)
	log.Printf("[EXPECT]: %v\n", expected)
}

// Test various patterns of arguments
// TODO: Add more patterns
var getTestArgs = [...][]interface{}{
	[]interface{}{nil, nil},
	[]interface{}{nil, nil, nil},
	[]interface{}{nil, authMap},
	[]interface{}{nil, authStruct},
	[]interface{}{nil, authMap, headerMap},
	[]interface{}{nil, authStruct, headerMap},
	[]interface{}{bodyMap, nil},
	[]interface{}{bodyMap, nil, nil},
	[]interface{}{bodyMap, nil, headerMap},
	[]interface{}{bodyMap, authMap},
	[]interface{}{bodyMap, authStruct},
	[]interface{}{bodyStruct, authMap},
	[]interface{}{bodyHybridMap, authMap},
	[]interface{}{bodyHybridMap, authStruct},
	[]interface{}{bodyMap, authMap, headerMap},
	[]interface{}{bodyMap, authStruct, headerMap},
	[]interface{}{bodyMap, authStruct, headerStruct},
	[]interface{}{bodyStruct, authMap, headerMap},
	[]interface{}{bodyStruct, authStruct, headerMap},
	[]interface{}{bodyStruct, authStruct, headerStruct},
	[]interface{}{bodyHybridMap, authMap, headerMap},
	[]interface{}{bodyHybridMap, authStruct, headerMap},
	[]interface{}{bodyHybridMap, authStruct, headerStruct},
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

func multTypeHandler(w http.ResponseWriter, r *http.Request) {
	data := []byte(`{"foo": ["bar", "baz"]}`)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

// argsHandler writes the body's key if the body exists and
// "auth" if Authorization key is founded in the requests' header.
func argsHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err == io.EOF {
		io.WriteString(w, "nada")
	}
	if len(data) > 0 {
		for key := range data {
			io.WriteString(w, key)
		}
	}
	if len(r.Header["Authorization"]) > 0 {
		io.WriteString(w, "authy")
	}
}

func TestGetResponseType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, args := range getTestArgs {
		resp, err := Get(ts.URL, args...)
		if err != nil {
			t.Error(err)
		}
		if reflect.TypeOf(resp) != reflect.TypeOf(&Response{}) {
			t.Error(err)
		}
	}
}

func TestGetResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	resp, err := Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "Hello world!" {
		log.Println(string(body))
		t.Error(err)
	}
}

func TestGetResponseAsBytes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	resp, err := Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	result := resp.Bytes()
	expected := []byte("Hello world!")
	if !reflect.DeepEqual(result, expected) {
		t.Error(err)
	}
}

var jsonTestTable = []struct {
	fn       func(http.ResponseWriter, *http.Request)
	expected []byte
}{
	{jsonHandler, []byte(`{"foo": ["bar", "baz"]}`)},
	{jsonWithTypeParamHandler, []byte(`{"foo": ["bar", "baz"]}`)},
	{multTypeHandler, []byte(`{"foo": ["bar", "baz"]}`)},
}

func TestGetResponseAsJSON(t *testing.T) {
	for _, tt := range jsonTestTable {
		ts := httptest.NewServer(http.HandlerFunc(tt.fn))
		defer ts.Close()
		resp, err := Get(ts.URL)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(resp.JSON(), tt.expected) {
			t.Error(err)
		}
	}
}

func TestGetResponseAsString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	resp, err := Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	result := resp.String()
	expected := "Hello world!"
	if result != expected {
		t.Error(err)
	}
}

var getTestTable = []struct {
	args     []interface{}
	expected string
}{
	{[]interface{}{}, "nada"},
	{[]interface{}{nil}, "nada"},
	{[]interface{}{bodyMap}, "foo"},
	{[]interface{}{bodyHybridMap}, "duplica"},
	{[]interface{}{nil, authMap}, "nadaauthy"},
	{[]interface{}{bodyMap, authMap}, "fooauthy"},
	{[]interface{}{bodyHybridMap, authMap}, "duplicaauthy"},
}

func TestGetArgsConcatInResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(argsHandler))
	defer ts.Close()
	for _, tt := range getTestTable {
		resp, err := Get(ts.URL, tt.args...)
		if err != nil {
			t.Error(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		if string(body) != tt.expected {
			log.Printf("[ACTUAL]: %v\n", string(body))
			log.Printf("[EXPECTED]: %v\n", tt.expected)
			t.Error(err)
		}
	}
}

var (
	badURLs = []string{
		"://maggot.#&",
		"crap://bs.com",
		"htp://f#as3",
	}
	badArgs = []interface{}{
		make(chan int),
		func() {},
		12i,
	}
)

func TestGetWithBadURLs(t *testing.T) {
	for _, url := range badURLs {
		_, err := Get(url)
		if err == nil {
			t.Error(err)
		}
	}
}

func TestGetWithBadData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, arg := range badArgs {
		_, err := Get(ts.URL, arg)
		if err == nil {
			t.Error(err)
		}
	}
}

func TestGetWithBadAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, arg := range badArgs {
		_, err := Get(ts.URL, nil, arg)
		if err == nil {
			t.Error(err)
		}
	}
}

func TestGetWithBadHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, arg := range badArgs {
		_, err := Get(ts.URL, nil, nil, arg)
		if err == nil {
			t.Error(err)
		}
	}
}

func TestGetWithForthArgs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	extra := map[string]string{"foo": "bar"}
	_, err := Get(ts.URL, nil, nil, nil, extra)
	if err == nil {
		t.Error(err)
	}
}

var getSyncTestTable = []struct {
	delay    int
	expected int
}{
	{1, 2},
	{2, 4},
	{3, 6},
	{4, 8},
}

// Get should wait for the response and return
func TestGetResponseTimes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, tt := range getSyncTestTable {
		delay := time.Duration(tt.delay) * time.Second
		expected := time.Duration(tt.expected) * time.Second
		p := relay.NewProxy(delay, ts)
		start := time.Now()
		_, _ = Get(p.URL)
		elapsed := time.Since(start)
		if elapsed <= expected {
			logExpectedResult(elapsed, expected)
			t.Error("Client returned before it should.")
		}
	}
}

var timeoutTestTable = []struct {
	delay    int
	timeout  float64
	expected float64
}{
	{1, 0.5, 0.5},
	{2, 0.5, 0.5},
	{2, 1.0, 1.0},
	{3, 1.0, 1.0},
}

// Get should wait fo the response until timed out.
func TestGetResponseOnTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()
	for _, tt := range timeoutTestTable {
		delay := time.Duration(tt.delay) * time.Second
		p := relay.NewProxy(delay, ts)
		start := time.Now()
		_, err := POCGet(p.URL, nil, nil, nil, tt.timeout)
		if err == nil {
			t.Error(errors.New("Client did not time out."))
		}
		elapsed := time.Since(start).Seconds()
		deviation := floatTimeDev.Seconds()
		if !(elapsed >= tt.expected-deviation || elapsed <= tt.expected+deviation) {
			logExpectedResult(elapsed, tt.expected)
			t.Error("Client returned before it should.")
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
	for _, tt := range getSyncTestTable {
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
				for k, _ := range d {
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
				for k, _ := range d {
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
