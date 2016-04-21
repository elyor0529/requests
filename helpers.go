package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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

/*
func marshalAll(args ...interface{}) (map[string][]byte, error) {

	res := make(map[string][]byte)
	typeError := errors.New("Unexpected type!")

	for i, arg := range args {
		switch i {
		case 0:
			method, ok := arg.(string)
			if !ok {
				return nil, typeError
			}
			res["method"] = []byte(method)
		case 1:
			url, ok := arg.(string)
			if !ok {
				return nil, typeError
			}
			res["url"] = []byte(url)
		case 2:
			body, err := json.Marshal(arg)
			if err != nil {
				return nil, err
			}
			res["body"] = body
		case 3:
			auth, err := json.Marshal(arg)
			if err != nil {
				return nil, err
			}
			res["auth"] = auth
		case 4:
			header, err := json.Marshal(arg)
			if err != nil {
				return nil, err
			}
			res["header"] = header
		}
	}

	return res, nil
}
*/

func marshalGet(args []interface{}) (map[string][]byte, error) {
	res := make(map[string][]byte)
	for i, arg := range args {
		switch i {
		case 0:
			if arg != nil {
				body, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				res["body"] = body
			}
		case 1:
			if arg != nil {
				auth, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				res["auth"] = auth
			}
		case 2:
			if arg != nil {
				header, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				res["header"] = header
			}
		default:
			return nil, fmt.Errorf(
				"Mismatched number of arguments: %d > %d",
				len(args), 4)
		}
	}
	return res, nil
}

func marshalGetAll(args []interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	for i, arg := range args {
		switch i {
		case 0:
			if arg != nil {
				body, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				res["body"] = bytes.NewBuffer(body)
			}
		case 1:
			if arg != nil {
				auth, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				res["auth"] = bytes.NewBuffer(auth)
			}
		case 2:
			if arg != nil {
				headers, ok := arg.(map[string][]string)
				if !ok {
					panic(errors.New("Mismatched type: header. Expect `map[string][]string`"))
				}
				var header http.Header = headers
				res["header"] = header
			}
		case 3:
			if arg != nil {
				sec, ok := arg.(float64)
				if !ok {
					panic(errors.New("Mismatched type: timeout. Expect `float64`"))
				}
				res["timeout"] = Ftos(sec)
			}
		default:
			panic(errors.New("Mismatched number of arguments."))
		}
	}
	return res, nil
}
