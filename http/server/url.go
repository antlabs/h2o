package server

import "strings"

const (
	https = "https://"
	http  = "http"
)

func TakeURL(u string, isTemplate bool) string {
	if strings.HasPrefix(u, https) {
		u = u[len(https)+1:]
	} else if strings.HasPrefix(u, http) {
		u = u[len(http)+1:]
	}

	pos := strings.Index(u, "/")
	if pos != -1 {
		u = u[pos:]
	}
	if !isTemplate {
		return u
	}

	u = strings.ReplaceAll(u, "{{.", ":")
	return strings.ReplaceAll(u, "}}", "")
}
