package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Fantom-foundation/go-evm/src/engine"
	"github.com/Fantom-foundation/go-lachesis/src/utils"
)


//AddRunFlags adds flags to the Run command
func AddRunFlags(cmd *cobra.Command) {
	//Lachesis Socket
	cmd.Flags().String("proxy", config.ProxyAddr, "IP:PORT of Lachesis proxy")
	if runtime.GOOS != "windows" {
		cmd.Flags().String("pidfile", config.Pidfile, "pidfile location; /tmp/go-evm.pid by default")
	}
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

	if runtime.GOOS != "windows" {
		err := utils.CheckPid(config.Pidfile)
		if err != nil {
			return err
		}
	}
	engine, err := engine.NewSocketEngine(*config, logger)
	//engine, err := engine.NewInmemEngine(*config, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
