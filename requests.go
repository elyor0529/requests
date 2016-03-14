package requests

import (
	"fmt"
	"io"
	"net/http"
)

type Auth struct {
	user     string
	password string
}

func Get(url string, body io.Reader, auth Auth) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		panic(err)
	}
	/*
		if auth {
			req.SetBasicAuth(auth.user, auth.password)
		}
	*/
	client := &http.Client{
		Timeout: 2,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Bye")
	return resp, nil
}

func GetAsync(url string, body io.Reader, auth Auth) (*http.Response, error) {
	response := make(chan *http.Response)
	go func() {
		req, err := http.NewRequest("GET", url, body)
		if err != nil {
			panic(err)
		}
		client := &http.Client{
			CheckRedirect: nil,
			Timeout:       0,
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		response <- resp
	}()
	return <-response, nil
}

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
