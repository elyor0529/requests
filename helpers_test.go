package requests

import (
	"errors"
	"log"
	"reflect"
	"testing"
)

var (
	getURL  = "http://httpbin.org/get"
	postURL = "http://httpbin.org/post"
	body    = map[string]string{"foo": "bar"}
	auth    = map[string]string{"user": "pass"}
	header  = map[string][]string{"Accept": {"text/plain"}}
)

var marshalGetTestTable = []struct {
	args     []interface{}
	expected map[string][]byte
}{
	{[]interface{}{}, map[string][]byte{}},
	{[]interface{}{nil}, map[string][]byte{}},
	{[]interface{}{nil, nil}, map[string][]byte{}},
	{[]interface{}{nil, nil, nil}, map[string][]byte{}},

	{[]interface{}{body}, map[string][]byte{
		"body": []byte(`{"foo":"bar"}`),
	}},
	{[]interface{}{nil, auth}, map[string][]byte{
		"auth": []byte(`{"user":"pass"}`),
	}},
	{[]interface{}{body, auth}, map[string][]byte{
		"body": []byte(`{"foo":"bar"}`),
		"auth": []byte(`{"user":"pass"}`),
	}},
	{[]interface{}{body, auth, nil}, map[string][]byte{
		"body": []byte(`{"foo":"bar"}`),
		"auth": []byte(`{"user":"pass"}`),
	}},
	{[]interface{}{body, auth, header}, map[string][]byte{
		"body":   []byte(`{"foo":"bar"}`),
		"auth":   []byte(`{"user":"pass"}`),
		"header": []byte(`{"Accept":["text/plain"]}`),
	}},
}

func TestMarshalGet(t *testing.T) {

	e := errors.New("Unexpected result!")

	for _, tt := range marshalGetTestTable {
		result, err := marshalGet(tt.args)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(result, tt.expected) {
			log.Printf("Result is %v", result)
			log.Printf("Expected is %v", tt.expected)
			t.Error(e)
		}

	}
}

/*
var marshalAllTable1 = []struct {
	method   string
	expected map[string][]byte
}{
	{"GET", map[string][]byte{"method": []byte("GET")}},
	{"POST", map[string][]byte{"method": []byte("POST")}},
	{"PUT", map[string][]byte{"method": []byte("PUT")}},
	{"DELETE", map[string][]byte{"method": []byte("DELETE")}},
	{"HEAD", map[string][]byte{"method": []byte("HEAD")}},
	{"OPTIONS", map[string][]byte{"method": []byte("OPTIONS")}},
}

var marshalAllTable2 = []struct {
	args     []interface{}
	expected map[string][]byte
}{
	{
		[]interface{}{"GET", getURL},
		map[string][]byte{
			"method": []byte("GET"),
			"url":    []byte(getURL),
		},
	},
	{
		[]interface{}{"GET", getURL, bodyMap},
		map[string][]byte{
			"method": []byte("GET"),
			"url":    []byte(getURL),
			"body":   []byte(`{"foo":"bar"}`),
		},
	},
	{
		[]interface{}{"GET", getURL, bodyStruct},
		map[string][]byte{
			"method": []byte("GET"),
			"url":    []byte(getURL),
			"body":   []byte(`{"foo":"bar"}`),
		},
	},
	{
		[]interface{}{"POST", postURL, bodyMap},
		map[string][]byte{
			"method": []byte("POST"),
			"url":    []byte(postURL),
			"body":   []byte(`{"foo":"bar"}`),
		},
	},
}

func TestMarshalAll1(t *testing.T) {

	e := errors.New("Unexpected result!")

	for _, tt := range marshalAllTable1 {
		res, err := marshalAll(tt.method)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(res, tt.expected) {
			t.Error(e)
		}

	}
}

func TestMarshalAll2(t *testing.T) {

	e := errors.New("Unexpected result!")

	for _, tt := range marshalAllTable2 {
		res, err := marshalAll(tt.args...)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(res, tt.expected) {
			t.Logf("Res is %v\n", res)
			t.Log(string(res["method"]))
			t.Log(string(res["url"]))
			t.Log(string(res["body"]))
			//t.Error(e)
			t.Logf("Exp is %v\n", tt.expected)
			t.Log(string(tt.expected["method"]))
			t.Log(string(tt.expected["url"]))
			t.Log(string(tt.expected["body"]))
			t.Error(e)
		}
	}
}
*/
