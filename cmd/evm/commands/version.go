package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Fantom-foundation/evm/src/version"
)

// VersionCmd displays the version of the evm program being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
