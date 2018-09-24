package commands

import (
	"fmt"

	"github.com/andrecronje/evm/engine"
	"github.com/spf13/cobra"
)

//AddRunFlags adds flags to the Run command
func AddRunFlags(cmd *cobra.Command) {

	//Base
	cmd.Flags().String("datadir", config.BaseConfig.DataDir, "Top-level directory for configuration and data")
	cmd.Flags().String("log_level", config.BaseConfig.LogLevel, "debug, info, warn, error, fatal, panic")

	//Eth
	cmd.Flags().String("eth.genesis", config.Eth.Genesis, "Location of genesis file")
	cmd.Flags().String("eth.keystore", config.Eth.Keystore, "Location of Ethereum account keys")
	cmd.Flags().String("eth.pwd", config.Eth.PwdFile, "Password file to unlock accounts")
	cmd.Flags().String("eth.db", config.Eth.DbFile, "Eth database file")
	cmd.Flags().String("eth.api_addr", config.Eth.EthAPIAddr, "Address of HTTP API service")
	cmd.Flags().Int("eth.cache", config.Eth.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")

	//Lachesis Socket
	cmd.Flags().String("lachesis.proxy_addr", config.Lachesis.ProxyAddr, "IP:PORT of Lachesis proxy")
	cmd.Flags().String("lachesis.client_addr", config.Lachesis.ClientAddr, "IP:PORT to bind client proxy")

	//Lachesis Inmem
	cmd.Flags().String("lachesis.dir", config.Lachesis.Dir, "Directory contaning priv_key.pem and peers.json files")
	cmd.Flags().String("lachesis.node_addr", config.Lachesis.NodeAddr, "IP:PORT of Lachesis node")
	cmd.Flags().String("lachesis.api_addr", config.Lachesis.APIAddr, "IP:PORT of Lachesis HTTP API service")
	cmd.Flags().Int("lachesis.heartbeat", config.Lachesis.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Int("lachesis.tcp_timeout", config.Lachesis.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("lachesis.cache_size", config.Lachesis.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("lachesis.sync_limit", config.Lachesis.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("lachesis.max_pool", config.Lachesis.MaxPool, "Max number of pool connections")
	cmd.Flags().String("lachesis.store_type", config.Lachesis.StoreType, "badger,inmem")
	cmd.Flags().String("lachesis.store_path", config.Lachesis.StorePath, "File containing the store database")
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

	// engine, err := engine.NewSocketEngine(*config, logger)
	engine, err := engine.NewInmemEngine(*config, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
