package docs

import (
	"fmt"
	"os"
)

func GenerateMarkdown(path, title, body string) error {
	content := fmt.Sprintf("# %s\n\n%s\n", title, body)
	return os.WriteFile(path, []byte(content), 0o644)
}
