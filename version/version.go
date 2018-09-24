package evm

import (
	"strings"
)

const Maj = "0"
const Min = "3"
const Fix = "4"

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
