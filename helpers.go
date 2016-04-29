package requests

import "encoding/json"

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
