package performance

import (
	"encoding/json"
	"io"
)

func DecodeZeroCopy(r io.Reader, v any) error { return json.NewDecoder(r).Decode(v) }
func EncodeZeroCopy(w io.Writer, v any) error { return json.NewEncoder(w).Encode(v) }
