package lachesis

import (
	_lachesis "github.com/andrecronje/lachesis/src/lachesis"
	"github.com/andrecronje/evm/src/config"
	"github.com/andrecronje/evm/src/service"
	"github.com/andrecronje/evm/src/state"
	"github.com/sirupsen/logrus"
)

// InmemLachesis implementes the Consensus interface.
// It uses an inmemory Lachesis node.
type InmemLachesis struct {
	config       *config.LachesisConfig
	lachesis     *_lachesis.Lachesis
	ethService   *service.Service
	ethState     *state.State
	logger       *logrus.Logger
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
