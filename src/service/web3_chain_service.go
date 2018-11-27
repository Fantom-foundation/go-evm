package service

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type Web3ChainService struct {
	backend *Service
}

func NewWeb3ChainServiceConstructor(backend *Service) RpcServiceConstructor {
	return func(context *RpcServiceContext) (RpcService, error) {
		return &Web3ChainService{
			backend: backend,
		}, nil
	}
}

func (s *Web3ChainService) Start() error {
	return nil
}

func (s *Web3ChainService) Stop() error {
	return nil
}

func (s *Web3ChainService) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumChainAPI(s.backend),
			Public:    true,
		}, /*{
			Namespace: "eth",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		},*/{
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s.backend),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugChainAPI(s.backend),
			Public:    true,
		}, /*{
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugChainAPI(s.chainConfig, s),
		},*/ /*{
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},*/
	}
}
