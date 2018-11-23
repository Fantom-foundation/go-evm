package service

import (
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/Fantom-foundation/evm/src/service/internal/ethapi"
)

type EthService struct {
	backend ethapi.Backend
}

func NewEthServiceConstructor(backend ethapi.Backend) RpcServiceConstructor {
	return func(context *RpcServiceContext) (RpcService, error) {
		return &EthService{
			backend: backend,
		}, nil
	}
}

func (s *EthService) Start() error {
	return nil
}

func (s *EthService) Stop() error {
	return nil
}

func (s *EthService) APIs() []rpc.API {
	return ethapi.GetAPIs(s.backend)
}
