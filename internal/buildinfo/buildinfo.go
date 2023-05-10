// Package buildinfo exposes ldflag-injected build metadata (version, commit,
// date) to packages other than main.
//
// main calls Set once at startup with the values bound to main.version /
// main.commit / main.date (those names are what the Makefile and goreleaser
// config inject into). Everywhere else, read the package-level variables or
// call Display.
package buildinfo

import "strings"

const (
	defaultVersion = "dev"
	defaultCommit  = "unknown"
	defaultDate    = "unknown"
)

var (
	Version = defaultVersion
	Commit  = defaultCommit
	Date    = defaultDate
)

// Set updates build metadata. Empty values are ignored.
func Set(version, commit, date string) {
	if version != "" {
		Version = version
	}
	if commit != "" {
		Commit = commit
	}
	if date != "" {
		Date = date
	}
}

// Display returns a human-readable version string suitable for diagnostics
// and for the Telegram session DeviceConfig.AppVersion field. For tagged
// builds it returns the tag (e.g. "v0.2.0"); for dev builds it falls back
// to a short commit hash when one is available.
func Display() string {
	if Version != defaultVersion && Version != "" {
		return Version
	}
	if Commit != defaultCommit && Commit != "" {
		return "dev (" + shortCommit(Commit) + ")"
	}
	return defaultVersion
}

func shortCommit(commit string) string {
	commit = strings.TrimSpace(commit)
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}
