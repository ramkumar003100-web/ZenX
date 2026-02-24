package docs

import (
	"fmt"
	"strings"
)

func BuildChangelog(version string, entries []string) string {
	return fmt.Sprintf("## %s\n- %s\n", version, strings.Join(entries, "\n- "))
}
