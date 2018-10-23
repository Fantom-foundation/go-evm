package version

import (
	"strings"
)

const Maj = "0"
const Min = "4"
const Fix = "0"

var (
	// The full version string
	Version = strings.Join([]string{Maj, Min, Fix}, ".")

	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit[:8]
	}
}
