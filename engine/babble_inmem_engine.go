package engine

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/mosaicnetworks/babble/crypto"
	"github.com/mosaicnetworks/babble/hashgraph"
	"github.com/mosaicnetworks/babble/net"
	"github.com/mosaicnetworks/babble/node"
	"github.com/mosaicnetworks/evm-babble/service"
	"github.com/mosaicnetworks/evm-babble/state"
	"github.com/sirupsen/logrus"
)

type BabbleInmemEngine struct {
	service *service.Service
	state   *state.State
	node    *node.Node
}

func NewBabbleInmemEngine(config Config, logger *logrus.Logger) (*BabbleInmemEngine, error) {
	submitCh := make(chan []byte)

	state, err := state.NewState(logger,
		config.BaseConfig.DbFile,
		config.BaseConfig.Cache)
	if err != nil {
		return nil, err
	}

	service := service.NewService(config.BaseConfig.EthDir,
		config.BaseConfig.APIAddr,
		config.BaseConfig.PwdFile,
		state,
		submitCh,
		logger)

	appProxy := NewInmemProxy(state, service, submitCh, logger)

	//------------------------------------------------------------------------------

	// Create the PEM key
	pemKey := crypto.NewPemKey(config.Babble.BabbleDir)

	// Try a read
	key, err := pemKey.ReadKey()
	if err != nil {
		return nil, err
	}

	// Create the peer store
	peerStore := net.NewJSONPeers(config.Babble.PeersFile)
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
	}).Debug("PARTICIPANTS")

	conf := node.NewConfig(
		time.Duration(config.Babble.Heartbeat)*time.Millisecond,
		time.Duration(config.Babble.TCPTimeout)*time.Millisecond,
		config.Babble.CacheSize,
		config.Babble.SyncLimit,
		config.Babble.StoreType,
		config.Babble.StorePath,
		logger)

	//Instantiate the Store (inmem or badger)
	var store hashgraph.Store
	var needBootstrap bool
	switch conf.StoreType {
	case "inmem":
		store = hashgraph.NewInmemStore(pmap, conf.CacheSize)
	case "badger":
		//If the file already exists, load and bootstrap the store using the file
		if _, err := os.Stat(conf.StorePath); err == nil {
			logger.Debug("loading badger store from existing database")
			store, err = hashgraph.LoadBadgerStore(conf.CacheSize, conf.StorePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load BadgerStore from existing file: %s", err)
			}
			needBootstrap = true
		} else {
			//Otherwise create a new one
			logger.Debug("creating new badger store from fresh database")
			store, err = hashgraph.NewBadgerStore(pmap, conf.CacheSize, conf.StorePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create new BadgerStore: %s", err)
			}
		}
	default:
		return nil, fmt.Errorf("Invalid StoreType: %s", conf.StoreType)
	}

	trans, err := net.NewTCPTransport(
		config.Babble.BabbleAddr, nil, 2, conf.TCPTimeout, logger)
	if err != nil {
		return nil, fmt.Errorf("Creating TCP Transport: %s", err)
	}

	node := node.NewNode(conf, nodeID, key, peers, store, trans, appProxy)
	if err := node.Init(needBootstrap); err != nil {
		return nil, fmt.Errorf("Initializing node: %s", err)
	}

	return &BabbleInmemEngine{
		state:   state,
		service: service,
		node:    node,
	}, nil

}

/*******************************************************************************
Implement Engine interface
*******************************************************************************/

func (p *BabbleInmemEngine) Run() error {

	go p.service.Run()

	p.node.Run(true)

	return nil
}
