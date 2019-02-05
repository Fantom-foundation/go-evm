package commands

import (
	"fmt"

	"github.com/Fantom-foundation/go-evm/src/consensus/lachesis"
	"github.com/Fantom-foundation/go-evm/src/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//AddLachesisFlags adds flags to the Lachesis command
func AddLachesisFlags(cmd *cobra.Command) {
	cmd.Flags().String("lachesis.datadir", config.Lachesis.DataDir, "Directory contaning priv_key.pem and peers.json files")
	cmd.Flags().String("lachesis.listen", config.Lachesis.BindAddr, "IP:PORT of Lachesis node")
	cmd.Flags().String("lachesis.service-listen", config.Lachesis.ServiceAddr, "IP:PORT of Lachesis HTTP API service")
	cmd.Flags().Duration("lachesis.heartbeat", config.Lachesis.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("lachesis.timeout", config.Lachesis.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("lachesis.cache-size", config.Lachesis.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int64("lachesis.sync-limit", config.Lachesis.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("lachesis.max-pool", config.Lachesis.MaxPool, "Max number of pool connections")
	cmd.Flags().Bool("lachesis.store", config.Lachesis.Store, "use persistent store")
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}
}

//NewLachesisCmd returns the command that starts EVM-Lite with Lachesis consensus
func NewLachesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lachesis",
		Short: "Run the evm-lite node with Lachesis consensus",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			config.SetDataDir(config.BaseConfig.DataDir)

			logger.WithFields(logrus.Fields{
				"Lachesis": config.Lachesis,
			}).Debug("Config")

			return nil
		},
		RunE: runLachesis,
	}

	AddLachesisFlags(cmd)

	return cmd
}

func runLachesis(cmd *cobra.Command, args []string) error {

	lachesis := lachesis.NewInmemLachesis(config.Lachesis, logger)
	engine, err := engine.NewEngine(*config, lachesis, logger)
	if err != nil {
		return fmt.Errorf("error building Engine: %s", err)
	}

	if err := engine.Run(); err != nil {
		return fmt.Errorf("error running Engine: %s", err)
	}

	return nil
}
