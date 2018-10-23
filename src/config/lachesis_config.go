package config

import (
	"fmt"
	"time"

	_lachesis "github.com/andrecronje/lachesis/src/lachesis"
	"github.com/sirupsen/logrus"
)

var (
	defaultNodeAddr        = ":1337"
	defaultLachesisAPIAddr = ":8000"
	defaultHeartbeat       = 500 * time.Millisecond
	defaultTCPTimeout      = 1000 * time.Millisecond
	defaultCacheSize       = 50000
	defaultSyncLimit       = 1000
	defaultMaxPool         = 2
	defaultLachesisDir     = fmt.Sprintf("%s/lachesis", DefaultDataDir)
	defaultPeersFile       = fmt.Sprintf("%s/peers.json", defaultLachesisDir)
)

// LachesisConfig contains the configuration of a Lachesis node
type LachesisConfig struct {

	// Directory containing priv_key.pem and peers.json files
	DataDir string `mapstructure:"datadir"`

	// Address of Lachesis node (where it talks to other Lachesis nodes)
	BindAddr string `mapstructure:"listen"`

	// Lachesis HTTP API address
	ServiceAddr string `mapstructure:"service-listen"`

	// Gossip heartbeat
	Heartbeat time.Duration `mapstructure:"heartbeat"`

	// TCP timeout
	TCPTimeout time.Duration `mapstructure:"timeout"`

	// Max number of items in caches
	CacheSize int `mapstructure:"cache-size"`

	// Max number of Event in SyncResponse
	SyncLimit int `mapstructure:"sync-limit"`

	// Max number of connections in net pool
	MaxPool int `mapstructure:"max-pool"`

	// Database type; badger or inmeum
	Store bool `mapstructure:"store"`
}

// DefaultLachesisConfig returns the default configuration for a Lachesis node
func DefaultLachesisConfig() *LachesisConfig {
	return &LachesisConfig{
		DataDir:     defaultLachesisDir,
		BindAddr:    defaultNodeAddr,
		ServiceAddr: defaultLachesisAPIAddr,
		Heartbeat:   defaultHeartbeat,
		TCPTimeout:  defaultTCPTimeout,
		CacheSize:   defaultCacheSize,
		SyncLimit:   defaultSyncLimit,
		MaxPool:     defaultMaxPool,
	}
}

// SetDataDir updates the lachesis configuration directories if they were set to
// to default values.
func (c *LachesisConfig) SetDataDir(datadir string) {
	if c.DataDir == defaultLachesisDir {
		c.DataDir = datadir
	}
}

// ToRealLachesisConfig converts an evm/src/config.LachesisConfig to a
// lachesis/src/lachesis.LachesisConfig as used by Lachesis
func (c *LachesisConfig) ToRealLachesisConfig(logger *logrus.Logger) *_lachesis.LachesisConfig {
	lachesisConfig := _lachesis.NewDefaultConfig()
	lachesisConfig.DataDir = c.DataDir
	lachesisConfig.BindAddr = c.BindAddr
	lachesisConfig.ServiceAddr = c.ServiceAddr
	lachesisConfig.MaxPool = c.MaxPool
	lachesisConfig.Store = c.Store
	lachesisConfig.Logger = logger
	lachesisConfig.NodeConfig.HeartbeatTimeout = c.Heartbeat
	lachesisConfig.NodeConfig.TCPTimeout = c.TCPTimeout
	lachesisConfig.NodeConfig.CacheSize = c.CacheSize
	lachesisConfig.NodeConfig.SyncLimit = c.SyncLimit
	lachesisConfig.NodeConfig.Logger = logger
	return lachesisConfig
}
