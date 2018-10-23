package engine

import (
	"time"

	proxy "github.com/andrecronje/lachesis/src/proxy/socket/lachesis"
	"github.com/andrecronje/lachesis/src/poset"
	"github.com/andrecronje/evm/src/service"
	"github.com/andrecronje/evm/src/state"
	"github.com/andrecronje/evm/src/config"
	"github.com/sirupsen/logrus"
)

type SocketEngine struct {
	service  *service.Service
	state    *state.State
	proxy    *proxy.SocketLachesisProxy
	submitCh chan []byte
	logger   *logrus.Logger
}

func NewSocketEngine(config config.Config, logger *logrus.Logger) (*SocketEngine, error) {
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

	lproxy, err := proxy.NewSocketLachesisProxy(config.ProxyAddr,
		config.ClientAddr,
		NewHandler(state),
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
		/*case commit := <-s.proxy.CommitCh():
			s.logger.Debug("CommitBlock")
			stateHash, err := s.state.ProcessBlock(commit.Block)
			commit.Respond(stateHash.Bytes(), err)*/
		}
	}
}

// Implements proxy.ProxyHandler interface
type Handler struct {
      stateHash []byte
			state     *state.State
}

// Called when a new block is comming. This particular example just computes
// the stateHash incrementaly with incoming blocks
func (h *Handler) CommitHandler(block poset.Block) (stateHash []byte, err error) {
      /*hash := h.stateHash

      for _, tx := range block.Transactions() {
              hash = crypto.SimpleHashFromTwoHashes(hash, crypto.SHA256(tx))
      }

      h.stateHash = hash

      return h.stateHash, nil*/
			hash, err := h.state.ProcessBlock(block)
			return hash.Bytes(), nil
}

// Called when syncing with the network
func (h *Handler) SnapshotHandler(blockIndex int) (snapshot []byte, err error) {
      return []byte{}, nil
}

// Called when syncing with the network
func (h *Handler) RestoreHandler(snapshot []byte) (stateHash []byte, err error) {
      return []byte{}, nil
}

func NewHandler(state *state.State) *Handler {
      return &Handler{
				state:    state,
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
