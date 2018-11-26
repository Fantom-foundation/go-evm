package main

import (
	cmd "github.com/Fantom-foundation/go-evm/cmd/evm/commands"
)

func main() {

	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.NewSoloCmd(),
		cmd.NewRaftCmd(),
		cmd.NewRunCmd(),
		cmd.VersionCmd)

	//Do not print usage when error occurs
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
