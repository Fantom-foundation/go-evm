package lachesis

import (
	"github.com/sirupsen/logrus"

	"github.com/Fantom-foundation/go-evm/src/config"
	"github.com/Fantom-foundation/go-evm/src/service"
	"github.com/Fantom-foundation/go-evm/src/state"
	_lachesis "github.com/Fantom-foundation/go-lachesis/src/lachesis"
)

// InmemLachesis implementes the Consensus interface.
// It uses an inmemory Lachesis node.
type InmemLachesis struct {
	config     *config.LachesisConfig
	lachesis   *_lachesis.Lachesis
	ethService *service.Service
	ethState   *state.State
	logger     *logrus.Logger
}

// NewInmemLachesis instantiates a new InmemLachesis consensus system
func NewInmemLachesis(config *config.LachesisConfig, logger *logrus.Logger) *InmemLachesis {
	return &InmemLachesis{
		config: config,
		logger: logger,
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

// Init instantiates a Lachesis inmemory node
func (b *InmemLachesis) Init(state *state.State, service *service.Service) error {

	b.logger.Debug("INIT")

	b.ethState = state
	b.ethService = service

	realConfig := b.config.ToRealLachesisConfig(b.logger)
	realConfig.Proxy = NewInmemProxy(state, service, service.GetSubmitCh(), b.logger)

	lachesis := _lachesis.NewLachesis(realConfig)
	err := lachesis.Init()
	if err != nil {
		return err
	}
	b.lachesis = lachesis

	return nil
}

// Run starts the Lachesis node
func (b *InmemLachesis) Run() error {
	b.lachesis.Run()
	return nil
}

// Info returns Lachesis stats
func (b *InmemLachesis) Info() (map[string]string, error) {
	info := b.lachesis.Node.GetStats()
	info["type"] = "lachesis"
	return info, nil
}
