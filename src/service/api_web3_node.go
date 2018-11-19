package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"

	//"github.com/ethereum/go-ethereum/p2p/enode"
	//"github.com/ethereum/go-ethereum/metrics"

	"github.com/andrecronje/evm/src/config"
)

// GetNodeAPIs returns the collection of RPC descriptors this node offers.
func GetNodeAPIs(n *Service) []rpc.API {
	return []rpc.API{
		{
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewNodePrivateAdminAPI(n),
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewNodePublicAdminAPI(n),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewNodePublicDebugAPI(n),
			Public:    true,
		}, {
			Namespace: "web3",
			Version:   "1.0",
			Service:   NewNodePublicWeb3API(n),
			Public:    true,
		},
	}
}

// PrivateAdminAPI is the collection of administrative API methods exposed only
// over a secure RPC channel.
type NodePrivateAdminAPI struct {
	node *Service // Service interfaced by this API
}

// NewPrivateAdminAPI creates a new API definition for the private admin methods
// of the node itself.
func NewNodePrivateAdminAPI(node *Service) *NodePrivateAdminAPI {
	return &NodePrivateAdminAPI{node: node}
}

// AddPeer requests connecting to a remote node, and also maintaining the new
// connection at all times, even reconnecting if it is lost.
func (api *NodePrivateAdminAPI) AddPeer(url string) (bool, error) {
	/*
		// Make sure the server is running, fail otherwise
		server := api.node.Server()
		if server == nil {
			return false, ErrNodeStopped
		}
		// Try to add the url as a static peer and return
		node, err := enode.ParseV4(url)
		if err != nil {
			return false, fmt.Errorf("invalid enode: %v", err)
		}
		server.AddPeer(node)
		return true, nil
	*/
	return false, ErrNotImplemented
}

// RemovePeer disconnects from a remote node if the connection exists
func (api *NodePrivateAdminAPI) RemovePeer(url string) (bool, error) {
	/*
		// Make sure the server is running, fail otherwise
		server := api.node.Server()
		if server == nil {
			return false, ErrNodeStopped
		}
		// Try to remove the url as a static peer and return
		node, err := enode.ParseV4(url)
		if err != nil {
			return false, fmt.Errorf("invalid enode: %v", err)
		}
		server.RemovePeer(node)
		return true, nil
	*/
	return false, ErrNotImplemented
}

// AddTrustedPeer allows a remote node to always connect, even if slots are full
func (api *NodePrivateAdminAPI) AddTrustedPeer(url string) (bool, error) {
	/*
		// Make sure the server is running, fail otherwise
		server := api.node.Server()
		if server == nil {
			return false, ErrNodeStopped
		}
		node, err := enode.ParseV4(url)
		if err != nil {
			return false, fmt.Errorf("invalid enode: %v", err)
		}
		server.AddTrustedPeer(node)
		return true, nil
	*/
	return false, ErrNotImplemented
}

// RemoveTrustedPeer removes a remote node from the trusted peer set, but it
// does not disconnect it automatically.
func (api *NodePrivateAdminAPI) RemoveTrustedPeer(url string) (bool, error) {
	/*
		// Make sure the server is running, fail otherwise
		server := api.node.Server()
		if server == nil {
			return false, ErrNodeStopped
		}
		node, err := enode.ParseV4(url)
		if err != nil {
			return false, fmt.Errorf("invalid enode: %v", err)
		}
		server.RemoveTrustedPeer(node)
		return true, nil
	*/
	return false, ErrNotImplemented
}

// PeerEvents creates an RPC subscription which receives peer events from the
// node's p2p.Server
func (api *NodePrivateAdminAPI) PeerEvents(ctx context.Context) (*rpc.Subscription, error) {
	/*
		// Make sure the server is running, fail otherwise
		server := api.node.Server()
		if server == nil {
			return nil, ErrNodeStopped
		}

		// Create the subscription
		notifier, supported := rpc.NotifierFromContext(ctx)
		if !supported {
			return nil, rpc.ErrNotificationsUnsupported
		}
		rpcSub := notifier.CreateSubscription()

		go func() {
			events := make(chan *p2p.PeerEvent)
			sub := server.SubscribeEvents(events)
			defer sub.Unsubscribe()

			for {
				select {
				case event := <-events:
					notifier.Notify(rpcSub.ID, event)
				case <-sub.Err():
					return
				case <-rpcSub.Err():
					return
				case <-notifier.Closed():
					return
				}
			}
		}()

		return rpcSub, nil
	*/
	return nil, ErrNotImplemented
}

// StartRPC starts the HTTP RPC API server.
func (api *NodePrivateAdminAPI) StartRPC(host *string, port *int, cors *string, apis *string, vhosts *string) (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()

	if api.node.httpHandler != nil {
		return false, fmt.Errorf("HTTP RPC already running on %s", api.node.httpEndpoint)
	}

	if host == nil {
		h := config.DefaultHTTPHost
		if api.node.rpcConfig.HTTPHost != "" {
			h = api.node.rpcConfig.HTTPHost
		}
		host = &h
	}
	if port == nil {
		port = &api.node.rpcConfig.HTTPPort
	}

	allowedOrigins := api.node.rpcConfig.HTTPCors
	if cors != nil {
		allowedOrigins = nil
		for _, origin := range strings.Split(*cors, ",") {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
		}
	}

	allowedVHosts := api.node.rpcConfig.HTTPVirtualHosts
	if vhosts != nil {
		allowedVHosts = nil
		for _, vhost := range strings.Split(*host, ",") {
			allowedVHosts = append(allowedVHosts, strings.TrimSpace(vhost))
		}
	}

	modules := api.node.httpWhitelist
	if apis != nil {
		modules = nil
		for _, m := range strings.Split(*apis, ",") {
			modules = append(modules, strings.TrimSpace(m))
		}
	}

	if err := api.node.startHTTP(fmt.Sprintf("%s:%d", *host, *port), api.node.rpcAPIs, modules, allowedOrigins, allowedVHosts, api.node.rpcConfig.HTTPTimeouts); err != nil {
		return false, err
	}
	return true, nil
}

// StopRPC terminates an already running HTTP RPC API endpoint.
func (api *NodePrivateAdminAPI) StopRPC() (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()

	if api.node.httpHandler == nil {
		return false, fmt.Errorf("HTTP RPC not running")
	}
	api.node.stopHTTP()
	return true, nil
}

// StartWS starts the websocket RPC API server.
func (api *NodePrivateAdminAPI) StartWS(host *string, port *int, allowedOrigins *string, apis *string) (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()

	if api.node.wsHandler != nil {
		return false, fmt.Errorf("WebSocket RPC already running on %s", api.node.wsEndpoint)
	}

	if host == nil {
		h := config.DefaultWSHost
		if api.node.rpcConfig.WSHost != "" {
			h = api.node.rpcConfig.WSHost
		}
		host = &h
	}
	if port == nil {
		port = &api.node.rpcConfig.WSPort
	}

	origins := api.node.rpcConfig.WSOrigins
	if allowedOrigins != nil {
		origins = nil
		for _, origin := range strings.Split(*allowedOrigins, ",") {
			origins = append(origins, strings.TrimSpace(origin))
		}
	}

	modules := api.node.rpcConfig.WSModules
	if apis != nil {
		modules = nil
		for _, m := range strings.Split(*apis, ",") {
			modules = append(modules, strings.TrimSpace(m))
		}
	}

	if err := api.node.startWS(fmt.Sprintf("%s:%d", *host, *port), api.node.rpcAPIs, modules, origins, api.node.rpcConfig.WSExposeAll); err != nil {
		return false, err
	}
	return true, nil
}

// StopWS terminates an already running websocket RPC API endpoint.
func (api *NodePrivateAdminAPI) StopWS() (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()

	if api.node.wsHandler == nil {
		return false, fmt.Errorf("WebSocket RPC not running")
	}
	api.node.stopWS()
	return true, nil
}

// PublicAdminAPI is the collection of administrative API methods exposed over
// both secure and unsecure RPC channels.
type NodePublicAdminAPI struct {
	node *Service // Node interfaced by this API
}

// NewPublicAdminAPI creates a new API definition for the public admin methods
// of the node itself.
func NewNodePublicAdminAPI(node *Service) *NodePublicAdminAPI {
	return &NodePublicAdminAPI{node: node}
}

// Peers retrieves all the information we know about each individual peer at the
// protocol granularity.
func (api *NodePublicAdminAPI) Peers() ([]*p2p.PeerInfo, error) {
	/*
		server := api.node.Server()
		if server == nil {
			return nil, ErrNodeStopped
		}
		return server.PeersInfo(), nil
	*/
	return nil, ErrNotImplemented
}

// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.
func (api *NodePublicAdminAPI) NodeInfo() (*p2p.NodeInfo, error) {
	/*
		server := api.node.Server()
		if server == nil {
			return nil, ErrNodeStopped
		}
		return server.NodeInfo(), nil
	*/
	return nil, ErrNotImplemented
}

// Datadir retrieves the current data directory the node is using.
func (api *NodePublicAdminAPI) Datadir() string {
	//return api.node.DataDir()
	return ErrNotImplemented.Error()
}

// PublicDebugAPI is the collection of debugging related API methods exposed over
// both secure and unsecure RPC channels.
type NodePublicDebugAPI struct {
	node *Service // Node interfaced by this API
}

// NewPublicDebugAPI creates a new API definition for the public debug methods
// of the node itself.
func NewNodePublicDebugAPI(node *Service) *NodePublicDebugAPI {
	return &NodePublicDebugAPI{node: node}
}

// Metrics retrieves all the known system metric collected by the node.
func (api *NodePublicDebugAPI) Metrics(raw bool) (map[string]interface{}, error) {
	/*
		// Create a rate formatter
		units := []string{"", "K", "M", "G", "T", "E", "P"}
		round := func(value float64, prec int) string {
			unit := 0
			for value >= 1000 {
				unit, value, prec = unit+1, value/1000, 2
			}
			return fmt.Sprintf(fmt.Sprintf("%%.%df%s", prec, units[unit]), value)
		}
		format := func(total float64, rate float64) string {
			return fmt.Sprintf("%s (%s/s)", round(total, 0), round(rate, 2))
		}
		// Iterate over all the metrics, and just dump for now
		counters := make(map[string]interface{})
		metrics.DefaultRegistry.Each(func(name string, metric interface{}) {
			// Create or retrieve the counter hierarchy for this metric
			root, parts := counters, strings.Split(name, "/")
			for _, part := range parts[:len(parts)-1] {
				if _, ok := root[part]; !ok {
					root[part] = make(map[string]interface{})
				}
				root = root[part].(map[string]interface{})
			}
			name = parts[len(parts)-1]

			// Fill the counter with the metric details, formatting if requested
			if raw {
				switch metric := metric.(type) {
				case metrics.Counter:
					root[name] = map[string]interface{}{
						"Overall": float64(metric.Count()),
					}

				case metrics.Meter:
					root[name] = map[string]interface{}{
						"AvgRate01Min": metric.Rate1(),
						"AvgRate05Min": metric.Rate5(),
						"AvgRate15Min": metric.Rate15(),
						"MeanRate":     metric.RateMean(),
						"Overall":      float64(metric.Count()),
					}

				case metrics.Timer:
					root[name] = map[string]interface{}{
						"AvgRate01Min": metric.Rate1(),
						"AvgRate05Min": metric.Rate5(),
						"AvgRate15Min": metric.Rate15(),
						"MeanRate":     metric.RateMean(),
						"Overall":      float64(metric.Count()),
						"Percentiles": map[string]interface{}{
							"5":  metric.Percentile(0.05),
							"20": metric.Percentile(0.2),
							"50": metric.Percentile(0.5),
							"80": metric.Percentile(0.8),
							"95": metric.Percentile(0.95),
						},
					}

				case metrics.ResettingTimer:
					t := metric.Snapshot()
					ps := t.Percentiles([]float64{5, 20, 50, 80, 95})
					root[name] = map[string]interface{}{
						"Measurements": len(t.Values()),
						"Mean":         t.Mean(),
						"Percentiles": map[string]interface{}{
							"5":  ps[0],
							"20": ps[1],
							"50": ps[2],
							"80": ps[3],
							"95": ps[4],
						},
					}

				default:
					root[name] = "Unknown metric type"
				}
			} else {
				switch metric := metric.(type) {
				case metrics.Counter:
					root[name] = map[string]interface{}{
						"Overall": float64(metric.Count()),
					}

				case metrics.Meter:
					root[name] = map[string]interface{}{
						"Avg01Min": format(metric.Rate1()*60, metric.Rate1()),
						"Avg05Min": format(metric.Rate5()*300, metric.Rate5()),
						"Avg15Min": format(metric.Rate15()*900, metric.Rate15()),
						"Overall":  format(float64(metric.Count()), metric.RateMean()),
					}

				case metrics.Timer:
					root[name] = map[string]interface{}{
						"Avg01Min": format(metric.Rate1()*60, metric.Rate1()),
						"Avg05Min": format(metric.Rate5()*300, metric.Rate5()),
						"Avg15Min": format(metric.Rate15()*900, metric.Rate15()),
						"Overall":  format(float64(metric.Count()), metric.RateMean()),
						"Maximum":  time.Duration(metric.Max()).String(),
						"Minimum":  time.Duration(metric.Min()).String(),
						"Percentiles": map[string]interface{}{
							"5":  time.Duration(metric.Percentile(0.05)).String(),
							"20": time.Duration(metric.Percentile(0.2)).String(),
							"50": time.Duration(metric.Percentile(0.5)).String(),
							"80": time.Duration(metric.Percentile(0.8)).String(),
							"95": time.Duration(metric.Percentile(0.95)).String(),
						},
					}

				case metrics.ResettingTimer:
					t := metric.Snapshot()
					ps := t.Percentiles([]float64{5, 20, 50, 80, 95})
					root[name] = map[string]interface{}{
						"Measurements": len(t.Values()),
						"Mean":         time.Duration(t.Mean()).String(),
						"Percentiles": map[string]interface{}{
							"5":  time.Duration(ps[0]).String(),
							"20": time.Duration(ps[1]).String(),
							"50": time.Duration(ps[2]).String(),
							"80": time.Duration(ps[3]).String(),
							"95": time.Duration(ps[4]).String(),
						},
					}

				default:
					root[name] = "Unknown metric type"
				}
			}
		})
		return counters, nil
	*/
	return nil, ErrNotImplemented
}

// PublicWeb3API offers helper utils
type NodePublicWeb3API struct {
	stack *Service
}

// NewPublicWeb3API creates a new Web3Service instance
func NewNodePublicWeb3API(stack *Service) *NodePublicWeb3API {
	return &NodePublicWeb3API{stack}
}

// ClientVersion returns the node name
func (s *NodePublicWeb3API) ClientVersion() string {
	// return s.stack.Server().Name
	return ErrNotImplemented.Error()
}

// Sha3 applies the ethereum sha3 implementation on the input.
// It assumes the input is hex encoded.
func (s *NodePublicWeb3API) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
