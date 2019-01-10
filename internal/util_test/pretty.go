package util_test

import (
	"encoding/json"
	"fmt"
)

// PrettyFormat formats any objects for test debug use.
func PrettyFormat(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("orig=%v, err=%v", v, err)
	}
	return string(b)
}
