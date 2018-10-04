package engine

import (
	"github.com/andrecronje/evm/service"
	"github.com/andrecronje/lachesis/poset"
	"github.com/sirupsen/logrus"
)

//InmemProxy implements the AppProxy interface
type InmemProxy struct {
	service  *service.Service
	submitCh chan []byte
	logger   *logrus.Logger
}

//NewInmemProxy initializes and return a new InmemProxy
func NewInmemProxy(
	service *service.Service,
	submitCh chan []byte,
	logger *logrus.Logger) *InmemProxy {

	return &InmemProxy{
		service:  service,
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
	stateHash, err := i.service.ProcessBlock(block)
	return stateHash, err
}

//TODO - Implement these two functions
func (i *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

func (i *InmemProxy) Restore(snapshot []byte) error {
	return nil
}
