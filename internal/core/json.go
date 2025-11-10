package core

import (
	"encoding/json"
	"fmt"
)

func Stringify(v any) string {
	jsonBytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}

	return string(jsonBytes)
}
