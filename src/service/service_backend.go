package service

import (
	//"context"
	//"math/big"
	"github.com/ethereum/go-ethereum/accounts"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/core"
	//"github.com/ethereum/go-ethereum/core/state"
	//"github.com/ethereum/go-ethereum/core/types"
	//"github.com/ethereum/go-ethereum/core/vm"
	//"github.com/ethereum/go-ethereum/eth/downloader"
	//"github.com/ethereum/go-ethereum/ethdb"
	//"github.com/ethereum/go-ethereum/event"
	//"github.com/ethereum/go-ethereum/params"
	//"github.com/ethereum/go-ethereum/rpc"
	//"github.com/Fantom-foundation/go-lachesis/src/poset"
	// internal/ethapi.Backend interface implementation
)

/*
 *   General Ethereum API
 */
/*
func (s *Service) Downloader() *downloader.Downloader {
	return nil
}

func (s *Service) ProtocolVersion() int {
	// TODO: return valid value
	return 3
}

func (s *Service) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return nil, ErrNotImplemented
}

func (s *Service) ChainDb() ethdb.Database {
	return nil
}

func (s *Service) EventMux() *event.TypeMux {
	return nil
}
*/
func (s *Service) AccountManager() *accounts.Manager {
	return s.am
}

/*
 *	 BlockChain API
 */
/*
func (s *Service) SetHead(number uint64) {
	return
}

func (s *Service) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	return nil, ErrNotImplemented
}

func (s *Service) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	return nil, ErrNotImplemented
}

func (s *Service) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	return nil, nil, ErrNotImplemented
}

func (s *Service) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return nil, ErrNotImplemented
}

func (s *Service) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return nil, ErrNotImplemented
}

func (s *Service) GetTd(blockHash common.Hash) *big.Int {
	return nil
}

func (s *Service) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	return nil, nil, ErrNotImplemented
}

func (s *Service) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return nil
}

func (s *Service) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return nil
}

func (s *Service) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return nil
}
*/
/*
 *	 TxPool API
 */
/*
func (s *Service) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	s.logger.Debugf("SendTx %v", signedTx)
	return ErrNotImplemented
}

func (s *Service) GetPoolTransactions() (types.Transactions, error) {
	return nil, ErrNotImplemented
}

func (s *Service) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	return nil
}

func (s *Service) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return 0, ErrNotImplemented
}

func (s *Service) Stats() (pending int, queued int) {
	return 0, 0
}

func (s *Service) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return nil, nil
}

func (s *Service) SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription {
	return nil
}

func (s *Service) ChainConfig() *params.ChainConfig {
	// TODO: custom config
	return &params.ChainConfig{}
}

func (s *Service) CurrentBlock() *types.Block {
	i := s.state.GetBlockIndex()
	block, err := s.state.GetBlockById(i)
	if err != nil {
		s.logger.Error(err)
		return &types.Block{}
	}
	return ConvBlock(block)
}
*/
