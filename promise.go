package requests

import (
	"net/http"
	"sync"
)

type Promise struct {
	wg  sync.WaitGroup
	res *http.Response
	err error
}

func NewPromise(f func() (*http.Response, error)) *Promise {
	p := &Promise{}
	p.wg.Add(1)
	go func() {
		p.res, p.err = f()
		p.wg.Done()
	}()
	return p
}

func (p *Promise) Then(r func() *http.Response, e func() error) interface{} {
	result := make(chan interface{})
	go func() {
		p.wg.Wait()
		if p.err != nil {
			result <- e()
			return
		}
		result <- r()
	}()
	return <-result
}
