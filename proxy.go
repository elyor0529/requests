package requests

import (
	"net/http"
	"net/http/httptest"
	"time"
)

type Proxy struct {
	*httptest.Server
	Latency time.Duration
	Backend *httptest.Server
}

func NewUnstartedProxy(latency time.Duration, backend *httptest.Server) *Proxy {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-time.Tick(latency)
		backend.Config.Handler.ServeHTTP(w, r)
		<-time.Tick(latency)
	}))
	proxy := &Proxy{
		Server:  server,
		Latency: latency,
		Backend: backend,
	}
	return proxy
}

func NewProxy(latency time.Duration, backend *httptest.Server) *Proxy {
	proxy := NewUnstartedProxy(latency, backend)
	proxy.Start()
	return proxy
}
