package output

import (
	"encoding/json"
	"fmt"
)

func JSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func KV(key string, value any) {
	fmt.Printf("%s: %v\n", key, value)
}
