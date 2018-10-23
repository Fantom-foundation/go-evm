package commands

import (
	"fmt"

	"github.com/andrecronje/evm/src/engine"
	"github.com/spf13/cobra"

)

//AddRunFlags adds flags to the Run command
func AddRunFlags(cmd *cobra.Command) {
	//Lachesis Socket
	cmd.Flags().Bool("standalone", config.Standalone, "Do not create a proxy")
	cmd.Flags().String("lachesis.proxy_addr", config.ProxyAddr, "IP:PORT of Lachesis proxy")
	cmd.Flags().String("lachesis.client_addr", config.ClientAddr, "IP:PORT to bind client proxy")

	//Lachesis Inmem
	cmd.Flags().String("lachesis.listen", config.Lachesis.BindAddr, "IP:PORT of Lachesis node")
	cmd.Flags().String("lachesis.api_addr", config.Lachesis.ServiceAddr, "IP:PORT of Lachesis HTTP API service")
	cmd.Flags().Duration("lachesis.heartbeat", config.Lachesis.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("lachesis.timeout", config.Lachesis.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("lachesis.cache_size", config.Lachesis.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("lachesis.sync_limit", config.Lachesis.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("lachesis.max_pool", config.Lachesis.MaxPool, "Max number of pool connections")
	cmd.Flags().Bool("lachesis.store_type", config.Lachesis.Store, "badger,inmem")
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
