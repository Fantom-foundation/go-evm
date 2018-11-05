package commands

import (
	"fmt"

	"github.com/andrecronje/evm/src/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

)

//AddRunFlags adds flags to the Run command
func AddRunFlags(cmd *cobra.Command) {
	//Lachesis Socket
	cmd.Flags().String("proxy", config.ProxyAddr, "IP:PORT of Lachesis proxy")
	viper.BindPFlags(cmd.Flags())
}

// NewRunCmd returns the command that allows the CLI to start a node.
func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the evm node",
		RunE:  run,
	}

	AddRunFlags(cmd)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {

	engine, err := engine.NewSocketEngine(*config, logger)
	//engine, err := engine.NewInmemEngine(*config, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
