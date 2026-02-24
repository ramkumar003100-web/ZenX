package validation

import (
	"encoding/json"
	"io"
	"net/http"
)

func BindAndValidate(r *http.Request, dst any, v *Validator) error {
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return err
	}
	return v.Struct(dst)
}
