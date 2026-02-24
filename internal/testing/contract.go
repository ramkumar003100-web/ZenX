package testing

import (
	"encoding/json"
	"os"
)

type Contract struct {
	Name     string `json:"name"`
	Request  any    `json:"request"`
	Response any    `json:"response"`
}

func SaveContract(path string, c Contract) error {
	b, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(path, b, 0o644)
}
