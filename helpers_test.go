package requests

/*
var (
	getURL  = "http://httpbin.org/get"
	postURL = "http://httpbin.org/post"
)

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
