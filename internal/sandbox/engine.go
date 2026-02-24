package sandbox

import (
	"context"
	"fmt"
)

type RuntimePermissions struct {
	Module string
	Action string
}

type Engine struct {
	Policies *EnginePolicyStore
	Threats  *ReputationCache
}

func (e Engine) Authorize(ctx context.Context, module, action, ip string) error {
	_ = ctx
	p, ok := e.Policies.Get(module)
	if !ok {
		return fmt.Errorf("policy missing for module %s", module)
	}
	if t, ok := e.Threats.Get(ip); ok && t.Score > 0.85 {
		return fmt.Errorf("request blocked by threat score")
	}
	switch action {
	case "db":
		if !p.AllowDB {
			return fmt.Errorf("db access denied")
		}
	case "network":
		if !p.AllowNetwork {
			return fmt.Errorf("network access denied")
		}
	case "read":
		if !p.AllowFileRead {
			return fmt.Errorf("file read denied")
		}
	case "write":
		if !p.AllowFileWrite {
			return fmt.Errorf("file write denied")
		}
	}
	return nil
}
