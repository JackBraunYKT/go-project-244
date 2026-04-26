package formatters

import (
	"code/differ"
	"encoding/json"
)

func FormatJSON(diff []differ.DiffNode) (string, error) {
	data, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
