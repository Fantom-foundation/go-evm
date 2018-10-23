package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	// Base
	defaultLogLevel = "debug"
	DefaultDataDir  = defaultHomeDir()
)

// Config contains de configuration for an EVM-Lite node
type Config struct {

	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for EVM and State
	Eth *EthConfig `mapstructure:"eth"`

	// Options for Lachesis consensus
	Lachesis *LachesisConfig `mapstructure:"lachesis"`

	// Options for Raft consensus
	Raft *RaftConfig `mapstructure:"raft"`

	ProxyAddr  string                  `mapstructure:"proxy-listen"`
	ClientAddr string                  `mapstructure:"client-connect"`
	Standalone bool                    `mapstructure:"standalone"`
}

// DefaultConfig returns the default configuration for an EVM-Lite node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig:   DefaultBaseConfig(),
		Eth:          DefaultEthConfig(),
		Lachesis:     DefaultLachesisConfig(),
		Raft:         DefaultRaftConfig(),
		ProxyAddr:    "127.0.0.1:1338",
		ClientAddr:   "127.0.0.1:1339",
	}
}

// SetDataDir updates the root data directory as well as the various lower config
// for eth and consensus
func (c *Config) SetDataDir(datadir string) {
	c.BaseConfig.DataDir = datadir
	if c.Eth != nil {
		c.Eth.SetDataDir(fmt.Sprintf("%s/eth", datadir))
	}
	if c.Lachesis != nil {
		c.Lachesis.SetDataDir(fmt.Sprintf("%s/lachesis", datadir))
	}
	if c.Raft != nil {
		c.Raft.SetDataDir(fmt.Sprintf("%s/raft", datadir))
	}
}

/*******************************************************************************
BASE CONFIG
*******************************************************************************/

// BaseConfig contains the top level configuration for an EVM-Lachesis node
type BaseConfig struct {

	// Top-level directory of evm-lachesis data
	DataDir string `mapstructure:"datadir"`

	// Debug, info, warn, error, fatal, panic
	LogLevel string `mapstructure:"log"`
}

// DefaultBaseConfig returns the default top-level configuration for EVM-Lachesis
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		DataDir:  DefaultDataDir,
		LogLevel: defaultLogLevel,
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
