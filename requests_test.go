package requests

import (
	//"fmt"
	"bytes"
	"net/http"
	"testing"
	//"errors"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"time"
	//"runtime"
)

// Stub `requests` namespace
type stubRequests struct {
	Get      func(string, string, map[string]string) (*stubResponse, error)
	GetAsync func(string, string, map[string]string, int) (chan *stubResponse, error)
	//Post     func(string, string, map[string]interface{}) (*http.Response, error)
}

// Stub connection
type stubConnection struct {
	channel chan interface{}
	done    chan bool
	latency time.Duration
}

func newStubConnection(latency time.Duration) *stubConnection {
	channel := make(chan interface{}, 1)
	done := make(chan bool, 1)
	connection := &stubConnection{
		channel: channel,
		done:    done,
		latency: latency,
	}
	return connection
}

func (c *stubConnection) Close() {
	close(c.channel)
	close(c.done)
}

// Stub `http.Request`
type stubRequest struct{ http.Request }

func newStubRequest(method, rawurl string, body io.ReadCloser) (*stubRequest, error) {
	uri, err := url.ParseRequestURI(rawurl)
	if err != nil {
		panic("Something's wrong with your URI")
	}
	request := &stubRequest{
		Request: http.Request{
			Method: method,
			URL:    uri,
			Body:   body,
		},
	}
	return request, nil
}

// Stub `http.Response`
type stubResponse struct{ http.Response }

func newStubResponse(status string, code int, header http.Header, body io.ReadCloser) *stubResponse {
	response := &stubResponse{
		Response: http.Response{
			Status:     status,
			StatusCode: code,
			Proto:      "HTTP/1.0",
			Header:     header,
			Body:       body,
		},
	}
	return response
}

// Stub `http.Server`
type stubServer struct {
	http.Server
	response *stubResponse
	latency  time.Duration
}

func newStubServer(addr string, res *stubResponse, lat time.Duration) *stubServer {
	server := &stubServer{
		Server:   http.Server{Addr: addr},
		response: res,
		latency:  lat,
	}
	return server
}

func (s *stubServer) Reply(code statusCode) *stubResponse {
	// Block for server's latency
	<-time.Tick(s.latency)
	//time.Sleep(s.latency)
	// Create status code and return the response
	s.response.StatusCode = (int)(code)
	// TODO: Assign Status
	return s.response
}

// Stub the `http.Client`
type stubClient struct{ http.Client }

func newStubClient(timeout time.Duration) *stubClient {
	client := &stubClient{
		Client: http.Client{Timeout: timeout},
	}
	return client
}

// Inject `*stubConnection` and `*stubServer` to simulate a server call
func (c *stubClient) Do(req *stubRequest, conn *stubConnection, server *stubServer) (*stubResponse, error) {
	// Block for the duration of `conn.latency` + `server.latency`
	// to simulaate real-world latencies and test timeoutn
	code := (statusCode)(server.response.StatusCode)

	conn.channel <- req

	go func(conn *stubConnection) {
		<-time.Tick(conn.latency)
		<-conn.channel
		// TODO: Do something with the receive request
		conn.channel <- server.Reply(code)
		<-time.Tick(server.latency)
		conn.done <- true
	}(conn)
	<-conn.done
	res := <-conn.channel
	return res.(*stubResponse), nil
}

var (
	// Setup connection
	networkLatency = time.Duration(100) * time.Millisecond
	conn           = newStubConnection(networkLatency)

	// Setup client
	timeout = time.Duration(3) * time.Second
	client  = newStubClient(timeout)

	// Setup server
	res           = newStubResponse("200 OK", 200, header, body)
	endpoint      = "http://jochasinga.io"
	serverLatency = time.Duration(100) * time.Millisecond
	server        = newStubServer(endpoint, res, serverLatency)

	// Setup request
	header   = http.Header{}
	jsonStr  = `{"foo": ["bar", "baz"]}`
	body     = ioutil.NopCloser(bytes.NewBuffer([]byte(jsonStr)))
	auth     = map[string]string{"user": "pass"}
	requests = &stubRequests{
		Get: func(url, body string, auth map[string]string) (*stubResponse, error) {
			// Convert body from string to io.ReadCloser
			bodyReadCloser := ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
			//req, err := http.NewRequest("GET", url, bodyReadCloser)
			req, err := newStubRequest("GET", url, bodyReadCloser)
			if err != nil {
				panic(err)
			}
			// TODO: include basic auth
			/*
				if len(auth) > 0 {
					for user, password := range auth {
						req.SetBasicAuth(user, password)
					}
				}
			*/
			res, err := client.Do(req, conn, server)
			if err != nil {
				panic(err)
			}
			return res, nil
		},
		GetAsync: func(url, body string, auth map[string]string, timeout int) (chan *stubResponse, error) {
			data := ioutil.NopCloser(bytes.NewBuffer([]byte(body)))

			req, err := newStubRequest("GET", url, data)
			if err != nil {
				panic(err)
			}

			temp := make(chan *stubResponse, 1)

			go func(t chan *stubResponse) {
				res, err := client.Do(req, conn, server)
				if err != nil {
					panic(err)
				}
				t <- res
			}(temp)
			return temp, nil
		},
	}
)

func TestGetResponseType(t *testing.T) {
	resp, err := requests.Get(endpoint, jsonStr, auth)
	if err != nil {
		t.Error(err)
	}
	returnType := reflect.TypeOf(resp)
	responseType := reflect.TypeOf((*stubResponse)(nil))
	if returnType != responseType {
		t.Errorf("Expected return type of `*stubResponse`, but it was %v instead.", returnType)
	}
}

func TestGetResponseStatus(t *testing.T) {
	resp, err := requests.Get(endpoint, jsonStr, auth)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected StatusCode `200`, but it was %s instead.", resp.Status)
	}
}

func TestGetAsyncResponseType(t *testing.T) {
	timeout := 1
	resultChan, err := requests.GetAsync(endpoint, jsonStr, auth, timeout)
	if err != nil {
		t.Error(err)
	}
	returnType := reflect.TypeOf(resultChan)
	responseType := reflect.TypeOf((chan *stubResponse)(nil))
	if returnType != responseType {
		t.Errorf("Expected return type of `chan *stubResponse`, but it was %v instead.", returnType)
	}
}

func TestGetAsyncResponseStatus(t *testing.T) {
	timeout := 1
	resultChan, err := requests.GetAsync(endpoint, jsonStr, auth, timeout)
	if err != nil {
		t.Error(err)
	}

	select {
	case result := <-resultChan:
		if result.StatusCode != 200 {
			t.Errorf("Expected Status of `200 OK`, but it was `%s` instead.", res.Status)
		}
		break
	// TODO: Fix this timeout
	case <-time.Tick(time.Duration(timeout) * time.Second):
		t.Log("time out!")
	}
}

func TestMain(m *testing.M) {
	v := m.Run()
	defer conn.Close()
	if v == 0 {
		os.Exit(1)
	}
	os.Exit(v)
}
