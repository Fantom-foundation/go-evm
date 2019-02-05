package config

import (
	"github.com/ethereum/go-ethereum/node"
)

const (
	DefaultHTTPHost = "127.0.0.1" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 8545      // Default TCP port for the HTTP RPC server
	DefaultWSHost   = "127.0.0.1" // Default host interface for the websocket RPC server
	DefaultWSPort   = 8546      // Default TCP port for the websocket RPC server
)

var (
	DefaultModules = []string{"admin", "personal", "txpool", "eth", "net", "web3", "miner", "debug"}
)

// DefaultRpcConfig contains reasonable default settings.
var DefaultRpcConfig *node.Config = &node.DefaultConfig
/*node.Config{
	DataDir:          DefaultDataDir,
	HTTPHost:         DefaultHTTPHost,
	HTTPPort:         DefaultHTTPPort,
	HTTPModules:      DefaultModules,
	HTTPVirtualHosts: []string{"*"},
	HTTPTimeouts:     rpc.DefaultHTTPTimeouts,
	WSHost:           DefaultWSHost,
	WSPort:           DefaultWSPort,
	WSModules:        DefaultModules,
}
*/
