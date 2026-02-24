package plugins

import "fmt"

type Descriptor struct {
	Name    string
	Version string
	Entry   string
}

func (d Descriptor) Validate() error {
	if d.Name == "" || d.Version == "" || d.Entry == "" {
		return fmt.Errorf("invalid plugin descriptor")
	}
	return nil
}
