package app

import "strings"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// VersionString returns a single-line version for --version output.
func VersionString() string {
	parts := make([]string, 0, 3)
	v := strings.TrimSpace(version)
	if v == "" {
		v = "dev"
	}
	parts = append(parts, v)

	if c := strings.TrimSpace(commit); c != "" && c != "none" {
		parts = append(parts, c)
	}

	if d := strings.TrimSpace(date); d != "" && d != "unknown" {
		parts = append(parts, d)
	}

	return strings.Join(parts, " ")
}
