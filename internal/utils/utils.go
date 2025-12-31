package utils

import "encoding/json"

func Copy(src, dest any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}
