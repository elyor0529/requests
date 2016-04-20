package requests

import (
	"encoding/json"
	"errors"
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
			return nil, errors.New("Expected 3 arguments: Got more.")
		}
	}

	return res, nil
}
