package testing

import "testing"

func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}
}
