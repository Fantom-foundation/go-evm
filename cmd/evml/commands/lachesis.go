package commands

import (
	"fmt"

	"github.com/Fantom-foundation/go-evm/src/consensus/lachesis"
	"github.com/Fantom-foundation/go-evm/src/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//AddBabbleFlags adds flags to the Babble command
func AddBabbleFlags(cmd *cobra.Command) {
	cmd.Flags().String("lachesis.datadir", config.Babble.DataDir, "Directory contaning priv_key.pem and peers.json files")
	cmd.Flags().String("lachesis.listen", config.Babble.BindAddr, "IP:PORT of Babble node")
	cmd.Flags().String("lachesis.service-listen", config.Babble.ServiceAddr, "IP:PORT of Babble HTTP API service")
	cmd.Flags().Duration("lachesis.heartbeat", config.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("lachesis.timeout", config.Babble.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("lachesis.cache-size", config.Babble.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("lachesis.sync-limit", config.Babble.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("lachesis.max-pool", config.Babble.MaxPool, "Max number of pool connections")
	cmd.Flags().Bool("lachesis.store", config.Babble.Store, "use persistent store")
	viper.BindPFlags(cmd.Flags())
}

//NewBabbleCmd returns the command that starts EVM-Lite with Babble consensus
func NewBabbleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lachesis",
		Short: "Run the evm-lite node with Babble consensus",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			config.SetDataDir(config.BaseConfig.DataDir)

			logger.WithFields(logrus.Fields{
				"Babble": config.Babble,
			}).Debug("Config")

			return nil
		},
		RunE: runBabble,
	}

	AddBabbleFlags(cmd)

	return cmd
}

func runBabble(cmd *cobra.Command, args []string) error {

	lachesis := lachesis.NewInmemBabble(config.Babble, logger)
	engine, err := engine.NewEngine(*config, lachesis, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
