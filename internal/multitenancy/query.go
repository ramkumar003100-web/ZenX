package multitenancy

import (
	"context"
	"fmt"
)

func ScopedQuery(ctx context.Context, base string) (string, error) {
	t, ok := TenantFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("tenant missing")
	}
	return fmt.Sprintf("%s /* tenant=%s */", base, t), nil
}
