package requests

import (
	"fmt"
	"bytes"
	"time"
	"net/http"
	"io/ioutil"
)

type Auth struct {
	user     string
	password string
}

func Get(url, body string, auth map[string]string) (*http.Response, error) {
	bodyReadCloser := ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	req, err := http.NewRequest("GET", url, bodyReadCloser)
	if err != nil {
		return (*http.Response)(nil), err
	}
	client := &http.Client{}
	// TODO: include basic auth
	fmt.Println("going to Do")
	res, err := client.Do(req)
	if err != nil {
		return (*http.Response)(nil), err
	}
	return res, nil
}

func GetAsync(url, body string, auth map[string]string, timeout int) (chan *http.Response, error) {
	bodyReadCloser := ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	req, err := http.NewRequest("GET", url, bodyReadCloser)
	if err != nil {
		return nil, err
	}
	client := &http.Client{ Timeout: time.Duration(timeout) * time.Second }
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
/*
func Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		CheckRedirect: nil,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
*/
