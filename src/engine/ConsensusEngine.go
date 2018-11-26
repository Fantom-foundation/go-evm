package engine

import (
	"github.com/Fantom-foundation/go-evm/src/config"
	"github.com/Fantom-foundation/go-evm/src/consensus"
	"github.com/Fantom-foundation/go-evm/src/service"
	"github.com/Fantom-foundation/go-evm/src/state"
	"github.com/sirupsen/logrus"
)

// ConsensusEngine is the actor that coordinates State, Service and Consensus
type ConsensusEngine struct {
	state     *state.State
	service   *service.Service
	consensus consensus.Consensus
}

// NewConsensusEngine instantiates a new ConsensusEngine with coupled State, Service, and Consensus
func NewConsensusEngine(config config.Config,
	consensus consensus.Consensus,
	logger *logrus.Logger) (*ConsensusEngine, error) {
	submitCh := make(chan []byte)

	state, err := state.NewState(logger,
		config.Eth.DbFile,
		config.Eth.Cache)
	if err != nil {
		return nil, err
	}

	service := service.NewService(config.Eth.Genesis,
		config.Eth.Keystore,
		config.Eth.EthAPIAddr,
		config.Eth.PwdFile,
		state,
		submitCh,
		logger)

	if err := consensus.Init(state, service); err != nil {
		return nil, err
	}

	service.SetInfoCallback(consensus.Info)

	engine := &ConsensusEngine{
		state:     state,
		service:   service,
		consensus: consensus,
	}

	return engine, nil
}

// Run starts the engine's Service asynchronously and starts the Consensus system
// synchronously
func (e *ConsensusEngine) Run() error {

	go e.service.Run()

	e.consensus.Run()

	return nil
}
