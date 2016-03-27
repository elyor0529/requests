package requests

import (
	"net/http"
	"net/http/httptest"
	"time"
)

type proxy struct {
	*httptest.Server
	latency time.Duration
	backend *httptest.Server
}

func newUnstartedProxy(latency time.Duration, backend *httptest.Server) *proxy {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-time.Tick(latency)
		backend.Config.Handler.ServeHTTP(w, r)
		<-time.Tick(latency)
	}))
	proxy := &proxy{
		Server:  server,
		latency: latency,
		backend: backend,
	}
	return proxy
}

func newProxy(latency time.Duration, backend *httptest.Server) *proxy {
	proxy := newUnstartedProxy(latency, backend)
	proxy.Start()
	return proxy
}
