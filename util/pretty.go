package util

import "encoding/json"

func JsonPretty(v interface{}) string {
	bz, _ := json.MarshalIndent(v, "", "  ")
	return string(bz)
}
