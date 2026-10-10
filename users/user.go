package users

import "strings"

func CheckUsername(username string) bool {
	if len(username) < 6 || strings.Contains(username, "admin") {
		return false
	}
	return true
}
