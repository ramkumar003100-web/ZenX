package auth

func HasRole(roles []string, wanted ...string) bool {
	set := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		set[r] = struct{}{}
	}
	for _, want := range wanted {
		if _, ok := set[want]; ok {
			return true
		}
	}
	return false
}
