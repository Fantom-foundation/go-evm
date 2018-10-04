package engine

import (
	"time"

	"github.com/andrecronje/evm/service"
	proxy "github.com/andrecronje/lachesis/proxy/lachesis"
	"github.com/sirupsen/logrus"
)

type SocketEngine struct {
	service  *service.Service
	proxy    *proxy.SocketLachesisProxy
	submitCh chan []byte
	logger   *logrus.Logger
}

func NewSocketEngine(config Config, logger *logrus.Logger) (*SocketEngine, error) {
	submitCh := make(chan []byte)

	service := service.NewService(config.Eth.States,
		config.Eth.Keystore,
		config.Eth.EthAPIAddr,
		config.Eth.PwdFile,
		config.Eth.DbDir,
		config.Eth.Cache,
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
			stateHash, err := s.service.ProcessBlock(commit.Block)
			commit.Respond(stateHash, err)
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
