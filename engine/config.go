package engine

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	//Base
	defaultLogLevel = "debug"
	defaultDataDir  = defaultHomeDir()

	//Eth
	defaultEthAPIAddr   = ":8080"
	defaultCache        = 128
	defaultEthDir       = fmt.Sprintf("%s/eth", defaultDataDir)
	defaultKeystoreFile = fmt.Sprintf("%s/keystore", defaultEthDir)
	defaultGenesisFile  = fmt.Sprintf("%s/genesis.json", defaultEthDir)
	defaultPwdFile      = fmt.Sprintf("%s/pwd.txt", defaultEthDir)
	defaultDbFile       = fmt.Sprintf("%s/chaindata", defaultEthDir)

	//Lachesis
	defaultProxyAddr   = ":1339"
	defaultClientAddr  = ":1338"
	defaultNodeAddr    = ":1337"
	defaultAPIAddr     = ":8000"
	defaultHeartbeat   = 500
	defaultTCPTimeout  = 1000
	defaultCacheSize   = 50000
	defaultSyncLimit   = 1000
	defaultMaxPool     = 2
	defaultStoreType   = "badger"
	defaultDir         = fmt.Sprintf("%s/lachesis", defaultDataDir)
	defaultPeersFile   = fmt.Sprintf("%s/peers.json", defaultDir)
	defaultStorePath   = fmt.Sprintf("%s/db", defaultDir)
)

//Config contains the configuration for an EVM node
type Config struct {

	//Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	//Options for EVM and State
	Eth *EthConfig `mapstructure:"eth"`

	//Options for Lachesis
	Lachesis *LachesisConfig `mapstructure:"lachesis"`
}

//DefaultConfig returns the default configuration for an EVM node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig:   DefaultBaseConfig(),
		Eth:          DefaultEthConfig(),
		Lachesis:     DefaultLachesisConfig(),
	}
}

/*******************************************************************************
BASE CONFIG
*******************************************************************************/

//BaseConfig contains the top level configuration for an EVM node
type BaseConfig struct {

	//Top-level directory of evm data
	DataDir string `mapstructure:"datadir"`

	//Debug, info, warn, error, fatal, panic
	LogLevel string `mapstructure:"log_level"`
}

//DefaultBaseConfig returns the default top-level configuration for EVM
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		DataDir:  defaultDataDir,
		LogLevel: defaultLogLevel,
	}
}

/*******************************************************************************
ETH CONFIG
*******************************************************************************/

//EthConfig contains the configuration relative to the accounts, EVM, trie/db,
//and service API
type EthConfig struct {

	//Genesis file
	Genesis string `mapstructure:"genesis"`

	//Location of ethereum account keys
	Keystore string `mapstructure:"keystore"`

	//File containing passwords to unlock ethereum accounts
	PwdFile string `mapstructure:"pwd"`

	//File containing the levelDB database
	DbFile string `mapstructure:"db"`

	//Address of HTTP API Service
	EthAPIAddr string `mapstructure:"api_addr"`

	//Megabytes of memory allocated to internal caching (min 16MB / database forced)
	Cache int `mapstructure:"cache"`
}

//DefaultEthConfig return the default configuration for Eth services
func DefaultEthConfig() *EthConfig {
	return &EthConfig{
		Genesis:    defaultGenesisFile,
		Keystore:   defaultKeystoreFile,
		PwdFile:    defaultPwdFile,
		DbFile:     defaultDbFile,
		EthAPIAddr: defaultEthAPIAddr,
		Cache:      defaultCache,
	}
}

/*******************************************************************************
LACHESIS CONFIG           XXX this should probably be in Lachesis itself XXX
*******************************************************************************/

//LachesisConfig contains the configuration of a Lachesis node
type LachesisConfig struct {

	/*********************************************
	SOCKET
	*********************************************/

	//Address of Lachesis proxy
	ProxyAddr string `mapstructure:"proxy_addr"`

	//Address of Lachesis client proxy
	ClientAddr string `mapstructure:"client_addr"`

	/*********************************************
	Inmem
	*********************************************/

	//Directory containing priv_key.pem and peers.json files
	Dir string `mapstructure:"dir"`

	//Address of Lachesis node (where it talks to other Lachesis nodes)
	NodeAddr string `mapstructure:"node_addr"`

	//Lachesis HTTP API address
	APIAddr string `mapstructure:"api_addr"`

	//Gossip heartbeat in milliseconds
	Heartbeat int `mapstructure:"heartbeat"`

	//TCP timeout in milliseconds
	TCPTimeout int `mapstructure:"tcp_timeout"`

	//Max number of items in caches
	CacheSize int `mapstructure:"cache_size"`

	//Max number of Event in SyncResponse
	SyncLimit int `mapstructure:"sync_limit"`

	//Max number of connections in net pool
	MaxPool int `mapstructure:"max_pool"`

	//Database type; badger or inmeum
	StoreType string `mapstructure:"store_type"`

	//If StoreType = badger, location of database file
	StorePath string `mapstructure:"store_path"`
}

//DefaultLachesisConfig returns the default configuration for a Lachesis node
func DefaultLachesisConfig() *LachesisConfig {
	return &LachesisConfig{
		ProxyAddr:  defaultProxyAddr,
		ClientAddr: defaultClientAddr,
		Dir:        defaultDir,
		NodeAddr:   defaultNodeAddr,
		APIAddr:    defaultAPIAddr,
		Heartbeat:  defaultHeartbeat,
		TCPTimeout: defaultTCPTimeout,
		CacheSize:  defaultCacheSize,
		SyncLimit:  defaultSyncLimit,
		MaxPool:    defaultMaxPool,
		StoreType:  defaultStoreType,
		StorePath:  defaultStorePath,
	}
}

/*******************************************************************************
FILE HELPERS
*******************************************************************************/

func defaultHomeDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "EVM")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "EVM")
		} else {
			return filepath.Join(home, ".evm")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
