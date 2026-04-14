package version

import "fmt"

var (
	Version = "devel"
	Commit  = "unknown"
	Date    = "unknown"
)

func String() string {
	return fmt.Sprintf("gpc %s (commit: %s, built: %s)", Version, Commit, Date)
}
