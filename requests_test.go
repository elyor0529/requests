package requests

import (
	"fmt"
	//"errors"
	"encoding/json"
	//"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	//"io"
	"io/ioutil"
	//"net/url"
	//"os"
	"reflect"
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetResponseTypeAndContent(t *testing.T) {
	Convey("GIVEN the Server Handler", t, func() {
		var d map[string]interface{}
		rq := New()
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
		Convey("WITH data and auth maps", func() {
			auth := map[string]string{"user": "password"}
			data := map[string][]string{"foo": []string{"bar", "baz"}}
			resp, err := rq.Get(ts.URL, data, auth)
			if err != nil {
				t.Error(err)
			}
			Convey("EXPECT Get() to return *httpResponse", func() {
				returnType := reflect.TypeOf(resp)
				responseType := reflect.TypeOf((*http.Response)(nil))
				So(returnType, ShouldEqual, responseType)
			})
			Convey("EXPECT Get() to return correct content", func() {
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo, auth")
			})

		})
		Convey("WITH data and auth structs", func() {
			auth := struct {
				User     string
				Password string
			}{ "user", "password" }
			
			data := struct {
				Foo []string `json:"foo"`
			}{ []string{"bar", "baz"} }
			resp, err := rq.Get(ts.URL, data, auth)
			if err != nil {
				t.Error(err)
			}
			Convey("EXPECT Get() to return *httpResponse", func() {
				returnType := reflect.TypeOf(resp)
				responseType := reflect.TypeOf((*http.Response)(nil))
				So(returnType, ShouldEqual, responseType)

			})
			Convey("EXPECT Get() to return correct content", func() {
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo, auth")
			})
		})

		Convey("WITH data and auth as nil", func() {
			resp, err := rq.Get(ts.URL, nil, nil)
			if err != nil {
				t.Error(err)
			}
			Convey("EXPECT Get() to return type *httpResponse", func() {
				returnType := reflect.TypeOf(resp)
				responseType := reflect.TypeOf((*http.Response)(nil))
				So(returnType, ShouldEqual, responseType)
			})
			Convey("EXPECT Get() to return correct content", func() {
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests")
			})
		})
		Reset(func() {
			ts.Close()
		})
	})
}

func TestGetAsyncResponseTypeAndContent(t *testing.T) {
	Convey("GIVEN the Server Handler with delay proxy", t, func() {
		var d map[string]interface{}
		timeout := time.Duration(5) * time.Second
		rq := New()
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
		
		Convey("WITH data and auth maps", func() {
			auth := map[string]string{"user": "password"}
			data := map[string][]string{"foo": []string{"bar", "baz"}}
			rc, err := rq.GetAsync(ts.URL, data, auth, timeout)
			if err != nil {
				t.Error(err)
			}
			Convey("EXPECT GetAsync() to return chan *httpResponse", func() {
				returnType := reflect.TypeOf(rc)
				responseType := reflect.TypeOf((chan *http.Response)(nil))
				So(returnType, ShouldEqual, responseType)
			})
			Convey("EXPECT GetAsync() to return correct content", func() {
				resp := <-rc
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo, auth")
			})
		})
			
		Convey("WITH data and auth structs", func() {
			auth := struct {
				User     string
				Password string
			}{ "user", "password" }
			
			data := struct {
				Foo []string `json:"foo"`
			}{ []string{"bar", "baz"} }
			rc, err := rq.GetAsync(ts.URL, data, auth, timeout)
			if err != nil {
				t.Error(err)
			}
			
			Convey("EXPECT GetAsync() to return *chan httpResponse", func() {
				returnType := reflect.TypeOf(rc)
				responseType := reflect.TypeOf((chan *http.Response)(nil))
				So(returnType, ShouldEqual, responseType)
			})
			
			Convey("EXPECT GetAsync() to return correct content", func() {
				resp := <-rc
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests, foo, auth")
			})
		})
		Convey("WITH data and auth as nil", func() {
			rc, err := rq.GetAsync(ts.URL, nil, nil, timeout)
			if err != nil {
				t.Error(err)
			}
			
			Convey("EXPECT GetAsync() to return type chan *httpResponse", func() {
				returnType := reflect.TypeOf(rc)
				responseType := reflect.TypeOf((chan *http.Response)(nil))
				So(returnType, ShouldEqual, responseType)
			})
			Convey("EXPECT GetAsync() to return correct content", func() {
				resp := <-rc
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
				greeting := string(body)
				So(greeting, ShouldEqual, "Hello, requests")
			})
		})
		Reset(func() {
			ts.Close()
		})
	})
}


