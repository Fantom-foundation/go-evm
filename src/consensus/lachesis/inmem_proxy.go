package lachesis

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/andrecronje/lachesis/src/poset"
	"github.com/andrecronje/evm/src/service"
	"github.com/andrecronje/evm/src/state"
	"github.com/sirupsen/logrus"
)

// InmemProxy implements the Lachesis AppProxy interface
type InmemProxy struct {
	service  *service.Service
	state    *state.State
	submitCh chan []byte
	logger   *logrus.Entry
}

// NewInmemProxy initializes and return a new InmemProxy
func NewInmemProxy(state *state.State,
	service *service.Service,
	submitCh chan []byte,
	logger *logrus.Logger) *InmemProxy {

	return &InmemProxy{
		service:  service,
		state:    state,
		submitCh: submitCh,
		logger:   logger.WithField("module", "lachesis/proxy"),
	}
}

/*******************************************************************************
Implement Lachesis AppProxy Interface
*******************************************************************************/

// SubmitCh is the channel through which the Service sends transactions to the
// node.
func (i *InmemProxy) SubmitCh() chan []byte {
	return i.submitCh
}

// CommitBlock commits Block to the State and expects the resulting state hash
func (i *InmemProxy) CommitBlock(block poset.Block) ([]byte, error) {
	i.logger.Debug("CommitBlock")

	blockHashBytes, err := block.Hash()
	blockHash := common.BytesToHash(blockHashBytes)

	for x, tx := range block.Transactions() {
		if err := i.state.ApplyTransaction(tx, x, blockHash); err != nil {
			return []byte{}, err
		}
	}

	hash, err := i.state.Commit()
	if err != nil {
		return []byte{}, err
	}

	return hash.Bytes(), nil
}

//TODO - Implement these two functions
func (i *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

func (i *InmemProxy) Restore(snapshot []byte) error {
	return nil
}
