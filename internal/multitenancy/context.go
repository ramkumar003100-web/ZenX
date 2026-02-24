package multitenancy

import "context"

type tenantKey struct{}

func WithTenant(ctx context.Context, tenant string) context.Context {
	return context.WithValue(ctx, tenantKey{}, tenant)
}
func TenantFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(tenantKey{}).(string)
	return v, ok
}
