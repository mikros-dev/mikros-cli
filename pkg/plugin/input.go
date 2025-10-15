package plugin

import (
	"encoding/json"
	"fmt"
	"strings"
)

func inputToMap(in string) (map[string]interface{}, error) {
	var (
		out map[string]interface{}
		dec = json.NewDecoder(strings.NewReader(in))
	)

	dec.UseNumber()
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("%v: %w", in, err)
	}

	return out, nil
}
