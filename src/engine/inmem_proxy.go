package engine

import (
	"github.com/andrecronje/lachesis/src/poset"
	"github.com/andrecronje/evm/src/service"
	"github.com/andrecronje/evm/src/state"
	"github.com/sirupsen/logrus"
)

//InmemProxy implements the AppProxy interface
type InmemProxy struct {
	service  *service.Service
	state    *state.State
	submitCh chan []byte
	logger   *logrus.Logger
}

//NewInmemProxy initializes and return a new InmemProxy
func NewInmemProxy(state *state.State,
	service *service.Service,
	submitCh chan []byte,
	logger *logrus.Logger) *InmemProxy {

	return &InmemProxy{
		service:  service,
		state:    state,
		submitCh: submitCh,
		logger:   logger,
	}
}

/*******************************************************************************
Implement AppProxy Interface
*******************************************************************************/

//SubmitCh is the channel through which the Service sends transactions to the
//node.
func (i *InmemProxy) SubmitCh() chan []byte {
	return i.submitCh
}

//CommitBlock commits Block to the State and expects the resulting state hash
func (i *InmemProxy) CommitBlock(block poset.Block) ([]byte, error) {
	i.logger.Debug("CommitBlock")
	stateHash, err := i.state.ProcessBlock(block)
	return stateHash.Bytes(), err
}

//TODO - Implement these two functions
func (i *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

func (i *InmemProxy) Restore(snapshot []byte) error {
	return nil
}
