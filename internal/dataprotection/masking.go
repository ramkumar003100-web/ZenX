package dataprotection

import "strings"

func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) < 2 {
		return "***"
	}
	return parts[0][:1] + "***@" + parts[1]
}

func DetectPII(input string) bool {
	return strings.Contains(input, "@") || strings.Contains(input, "ssn") || strings.Contains(input, "passport")
}
