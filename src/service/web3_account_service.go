package service

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type Web3AccountService struct {
	backend *Service
}

func NewWeb3AccountServiceConstructor(backend *Service) RpcServiceConstructor {
	return func(context *RpcServiceContext) (RpcService, error) {
		return &Web3AccountService{
			backend: backend,
		}, nil
	}
}

func (s *Web3AccountService) Start() error {
	return nil
}

func (s *Web3AccountService) Stop() error {
	return nil
}

func (s *Web3AccountService) APIs() []rpc.API {
	nonceLock := new(AddrLocker)
	return []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumAPI(s.backend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(s.backend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicTransactionPoolAPI(s.backend, nonceLock),
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(s.backend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s.backend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.backend),
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(s.backend.AccountManager()),
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountAPI(s.backend, nonceLock),
			Public:    false,
		},
	}
}
