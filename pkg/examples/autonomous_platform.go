package examples

import (
	"context"
	"time"

	"zenx/internal/cluster"
	"zenx/internal/controlplane"
	"zenx/internal/dataprotection"
	"zenx/internal/deployment"
	"zenx/internal/traffic"
)

func LeaderElectionExample(ctx context.Context) string {
	m := cluster.NewMembership(5 * time.Second)
	m.Heartbeat(ctx, cluster.Node{ID: "node-a", Address: "10.0.0.1", IsHealthy: true})
	m.Heartbeat(ctx, cluster.Node{ID: "node-b", Address: "10.0.0.2", IsHealthy: true})
	e := cluster.NewLeaderElector(m)
	e.Run(ctx, 50*time.Millisecond)
	time.Sleep(80 * time.Millisecond)
	return e.Leader()
}

func DynamicConfigChangeExample(ctx context.Context, c *controlplane.Controller) controlplane.ConfigVersion {
	return c.Update(ctx, map[string]any{"rate_limit": 220, "feature_x": true}, map[string]string{"source": "ops"})
}

func DistributedLockExample(ctx context.Context, l *cluster.DistributedLock) bool {
	acquired := l.Acquire(ctx, "billing-close", "node-a", 2*time.Second)
	if acquired {
		_ = l.Release(ctx, "billing-close", "node-a")
	}
	return acquired
}

func FieldLevelEncryptionExample() (string, error) {
	key := []byte("12345678901234567890123456789012")
	return dataprotection.EncryptField([]byte("secret@email.com"), key)
}

func CanaryDeploymentExample(hash uint64) string {
	c := deployment.Canary{StableVersion: "v1", CanaryVersion: "v2", CanaryPercent: 0.15}
	return c.RouteByHash(hash)
}

func TrafficOverloadRecoveryExample(detector *traffic.OverloadDetector, bal *traffic.Balancer) bool {
	over := detector.Evaluate(92, 20000)
	if over {
		bal.Drain("node-overloaded", true)
	}
	return over
}
