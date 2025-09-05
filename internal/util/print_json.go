package util

import (
	"encoding/json"
	"fmt"
)

func PrintJson(label string, value any) {
	formatted, err := json.MarshalIndent(value, "", "    ")

	if err == nil {
		fmt.Println(label, "\n", string(formatted))
	}
}
