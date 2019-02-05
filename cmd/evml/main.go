package main

import (
	cmd "github.com/Fantom-foundation/go-evm/cmd/evml/commands"
)

func main() {

	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.NewSoloCmd(),
		cmd.NewLachesisCmd(),
		cmd.NewRaftCmd(),
		cmd.VersionCmd)

	//Do not print usage when error occurs
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
