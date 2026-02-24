package billing

func EnforceQuota(used, limit int64) bool { return used < limit }
