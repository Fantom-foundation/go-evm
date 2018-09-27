package main

import (
	cmd "github.com/andrecronje/evm/cmd/evm/commands"
)

func main() {

	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.NewRunCmd(),
		cmd.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
