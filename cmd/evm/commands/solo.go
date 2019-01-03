package commands

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Fantom-foundation/go-evm/src/consensus/solo"
	"github.com/Fantom-foundation/go-evm/src/engine"
)

var genesisTemplate = `
{
	"alloc": {
			"{{.}}": {
					"balance": "1000000000000000000000000000"
			}
	}
}
`

var genesisAddress string

//AddSoloFlags adds flags to the Solo command
func AddSoloFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&genesisAddress, "genesis", "", "create genesis file specifying pre-funded account with given address")
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic("Unable to bind viper flags")
	}
}

//NewSoloCmd returns the command that starts EVM-Lite with Solo consensus
func NewSoloCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solo",
		Short: "Run the evm node with Solo consensus (no consensus)",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			config.SetDataDir(config.BaseConfig.DataDir)

			logger.WithFields(logrus.Fields{
				"Eth":     config.Eth,
				"genesis": genesisAddress,
			}).Debug("Config")

			if cmd.Flags().Changed("genesis") {
				logger.Debug("Writing genesis file")
				if err := createGenesis(config.Eth.Genesis, genesisAddress); err != nil {
					return err
				}
			}

			return nil
		},
		RunE: runSolo,
	}
	AddSoloFlags(cmd)
	return cmd
}

func createGenesis(genesisFile, genesisAddr string) error {

	if _, err := os.Stat(genesisFile); err == nil {
		logger.WithError(err).Error("Genesis file already exists. Cannot overwrite.")
		return err
	}

	t := template.New("genesis")
	t, err := t.Parse(genesisTemplate) // parsing of template string
	if err != nil {
		logger.WithError(err).Error("Parsing genesis template")
		return err
	}

	genDir := filepath.Dir(genesisFile)
	if _, err := os.Stat(genDir); os.IsNotExist(err) {
		err = os.MkdirAll(genDir, 0755)
		if err != nil {
			logger.WithError(err).Error("Creating base directory of genesis file")
			return err
		}
	}

	f, err := os.OpenFile(genesisFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.WithError(err).Errorf("Creating file %s", genesisFile)
		return err
	}

	return t.Execute(f, genesisAddr)
}

func runSolo(cmd *cobra.Command, args []string) error {
	soloConsensus := solo.NewSolo(logger)
	soloEngine, err := engine.NewConsensusEngine(*config, soloConsensus, logger)
	if err != nil {
		return fmt.Errorf("error building Engine: %s", err)
	}

	if err := soloEngine.Run(); err != nil {
		return fmt.Errorf("error running Engine: %s", err)
	}

	return nil
}
