package consensus

import (
	"github.com/Fantom-foundation/evm/src/service"
	"github.com/Fantom-foundation/evm/src/state"
)

// Consensus is the interface that abstracts the consensus system
type Consensus interface {
	Init(*state.State, *service.Service) error
	Run() error
	Info() (map[string]string, error)
}
