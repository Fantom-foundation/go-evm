package engine

import (
	"time"

	proxy "github.com/andrecronje/lachesis/src/proxy/lachesis"
	"github.com/andrecronje/evm/service"
	"github.com/andrecronje/evm/state"
	"github.com/sirupsen/logrus"
)

type SocketEngine struct {
	service  *service.Service
	state    *state.State
	proxy    *proxy.SocketLachesisProxy
	submitCh chan []byte
	logger   *logrus.Logger
}

func NewSocketEngine(config Config, logger *logrus.Logger) (*SocketEngine, error) {
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

	lproxy, err := proxy.NewSocketLachesisProxy(config.Lachesis.ProxyAddr,
		config.Lachesis.ClientAddr,
		time.Duration(config.Lachesis.TCPTimeout)*time.Millisecond,
		logger)
	if err != nil {
		return nil, err
	}

	return &SocketEngine{
		service:  service,
		state:    state,
		proxy:    lproxy,
		submitCh: submitCh,
		logger:   logger,
	}, nil
}

func (s *SocketEngine) serve() {
	for {
		select {
		case tx := <-s.submitCh:
			s.logger.Debug("proxy about to submit tx")
			if err := s.proxy.SubmitTx(tx); err != nil {
				s.logger.WithError(err).Error("SubmitTx")
			}
			s.logger.Debug("proxy submitted tx")
		case commit := <-s.proxy.CommitCh():
			s.logger.Debug("CommitBlock")
			stateHash, err := s.state.ProcessBlock(commit.Block)
			commit.Respond(stateHash.Bytes(), err)
		}
	}
}

/*******************************************************************************
Implement Engine interface
*******************************************************************************/

func (s *SocketEngine) Run() error {

	go s.service.Run()

	s.serve()

	return nil
}
