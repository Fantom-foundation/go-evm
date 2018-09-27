package engine

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/andrecronje/lachesis/crypto"
	"github.com/andrecronje/lachesis/poset"
	"github.com/andrecronje/lachesis/net"
	"github.com/andrecronje/lachesis/node"
	serv "github.com/andrecronje/lachesis/service"
	"github.com/andrecronje/evm/service"
	"github.com/andrecronje/evm/state"
	"github.com/sirupsen/logrus"
)

type InmemEngine struct {
	ethService *service.Service
	ethState   *state.State
	node       *node.Node
	service    *serv.Service
}

func NewInmemEngine(config Config, logger *logrus.Logger) (*InmemEngine, error) {
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

	appProxy := NewInmemProxy(state, service, submitCh, logger)

	//------------------------------------------------------------------------------

	// Create the PEM key
	pemKey := crypto.NewPemKey(config.Lachesis.Dir)

	// Try a read
	key, err := pemKey.ReadKey()
	if err != nil {
		return nil, err
	}

	// Create the peer store
	peerStore := net.NewJSONPeers(config.Lachesis.Dir)
	// Try a read
	peers, err := peerStore.Peers()
	if err != nil {
		return nil, err
	}

	// There should be at least two peers
	if len(peers) < 2 {
		return nil, fmt.Errorf("Should define at least two peers")
	}

	sort.Sort(net.ByPubKey(peers))
	pmap := make(map[string]int)
	for i, p := range peers {
		pmap[p.PubKeyHex] = i
	}

	//Find the ID of this node
	nodePub := fmt.Sprintf("0x%X", crypto.FromECDSAPub(&key.PublicKey))
	nodeID := pmap[nodePub]

	logger.WithFields(logrus.Fields{
		"pmap": pmap,
		"id":   nodeID,
	}).Debug("Participants")

	conf := node.NewConfig(
		time.Duration(config.Lachesis.Heartbeat)*time.Millisecond,
		time.Duration(config.Lachesis.TCPTimeout)*time.Millisecond,
		config.Lachesis.CacheSize,
		config.Lachesis.SyncLimit,
		config.Lachesis.StoreType,
		config.Lachesis.StorePath,
		logger)

	//Instantiate the Store (inmem or badger)
	var store poset.Store
	var needBootstrap bool
	switch conf.StoreType {
	case "inmem":
		store = poset.NewInmemStore(pmap, conf.CacheSize)
	case "badger":
		//If the file already exists, load and bootstrap the store using the file
		if _, err := os.Stat(conf.StorePath); err == nil {
			logger.Debug("loading badger store from existing database")
			store, err = poset.LoadBadgerStore(conf.CacheSize, conf.StorePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load BadgerStore from existing file: %s", err)
			}
			needBootstrap = true
		} else {
			//Otherwise create a new one
			logger.Debug("creating new badger store from fresh database")
			store, err = poset.NewBadgerStore(pmap, conf.CacheSize, conf.StorePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create new BadgerStore: %s", err)
			}
		}
	default:
		return nil, fmt.Errorf("Invalid StoreType: %s", conf.StoreType)
	}

	trans, err := net.NewTCPTransport(
		config.Lachesis.NodeAddr, nil, 2, conf.TCPTimeout, logger)
	if err != nil {
		return nil, fmt.Errorf("Creating TCP Transport: %s", err)
	}

	node := node.NewNode(conf, nodeID, key, peers, store, trans, appProxy)
	if err := node.Init(needBootstrap); err != nil {
		return nil, fmt.Errorf("Initializing node: %s", err)
	}

	lserv := serv.NewService(config.Lachesis.APIAddr, node, logger)

	return &InmemEngine{
		ethState:   state,
		ethService: service,
		node:       node,
		service:    lserv,
	}, nil

}

/*******************************************************************************
Implement Engine interface
*******************************************************************************/

func (i *InmemEngine) Run() error {

	//ETH API service
	go i.ethService.Run()

	//Lachesis API service
	go i.service.Serve()

	i.node.Run(true)

	return nil
}
