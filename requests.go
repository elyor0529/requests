package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Req interface {
	Get(string, interface{}, interface{}) (*http.Response, error)
	GetAsync(string, interface{}, interface{}, time.Duration) (chan *http.Response, error)
	Post(string, string, interface{}) (*http.Response, error)
}

type Requests struct {}

func New() *Requests {
	return new(Requests)
}

func marshalData(data, auth interface{}) (map[string][]byte, error) {
	
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	a, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	results := map[string][]byte{"data": d, "auth": a}
	return results, nil
}

func (r *Requests) Get(url string, data, auth interface{}) (*http.Response, error) {

	results, err := marshalData(data, auth)
	if err != nil {
		return (*http.Response)(nil), err
	}

	dat, aut := results["data"], results["auth"]

	dataReadCloser := ioutil.NopCloser(bytes.NewBuffer(dat))

	req, err := http.NewRequest("GET", url, dataReadCloser)
	if err != nil {
		return (*http.Response)(nil), err
	}

	var authData map[string]interface{}

	if err := json.Unmarshal(aut, &authData); err != nil {
		return (*http.Response)(nil), err
	}

	for user, password := range authData {
		pw, ok := password.(string)
		if !ok {
			return (*http.Response)(nil), err
		}
		req.SetBasicAuth(user, pw)
	}                                                                                         
 
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return (*http.Response)(nil), err
	}
	
	return res, nil
}

func (r *Requests) GetAsync(url string, data, auth interface{}, timeout time.Duration) (chan *http.Response, error) {
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

	if err := json.Unmarshal(aut, &authData); err != nil {
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

func (r *Requests) Post(url, bodyType string, data interface{}) (*http.Response, error) {
	
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

