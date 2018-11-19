package service

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
	//"github.com/andrecronje/evm/src/service/internal/ethapi"
)

// StartRPC starts the various API endpoints, terminating all in case of errors
func (n *Service) StartRPC() error {
	apis := n.apis()

	if err := n.startInProc(apis); err != nil {
		return err
	}
	if err := n.startIPC(apis); err != nil {
		n.stopInProc()
		return err
	}
	if err := n.startHTTP(n.httpEndpoint, apis, n.rpcConfig.HTTPModules, n.rpcConfig.HTTPCors, n.rpcConfig.HTTPVirtualHosts, n.rpcConfig.HTTPTimeouts); err != nil {
		n.stopIPC()
		n.stopInProc()
		return err
	}
	if err := n.startWS(n.wsEndpoint, apis, n.rpcConfig.WSModules, n.rpcConfig.WSOrigins, n.rpcConfig.WSExposeAll); err != nil {
		n.stopHTTP()
		n.stopIPC()
		n.stopInProc()
		return err
	}
	// All API endpoints started successfully
	n.rpcAPIs = apis
	n.stop = make(chan struct{})
	return nil
}

// StopRPC terminates a running node along with all it's services. In the node was
// not started, an error is returned.
func (n *Service) StopRPC() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// Terminate the API, services and the p2p server.
	n.stopWS()
	n.stopHTTP()
	n.stopIPC()
	n.rpcAPIs = nil

	// unblock n.Wait
	close(n.stop)

	return nil
}

// WaitUntilRPC blocks the thread until the Web3 API is stopped. If the it is not running
// at the time of invocation, the method immediately returns.
func (n *Service) WaitUntilRPC() {
	n.lock.RLock()
	stop := n.stop
	n.lock.RUnlock()

	<-stop
}

// RestartWeb3API terminates a running node and boots up a new one in its place. If the
// node isn't running, an error is returned.
func (n *Service) RestartWeb3API() error {
	if err := n.StopRPC(); err != nil {
		return err
	}
	if err := n.StartRPC(); err != nil {
		return err
	}
	return nil
}

// startInProc initializes an in-process RPC endpoint.
func (n *Service) startInProc(apis []rpc.API) error {
	// Register all the APIs exposed by the services
	handler := rpc.NewServer()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
		n.logger.Debug("InProc registered", "service", api.Service, "namespace", api.Namespace)
	}
	n.inprocHandler = handler
	return nil
}

// stopInProc terminates the in-process RPC endpoint.
func (n *Service) stopInProc() {
	if n.inprocHandler != nil {
		n.inprocHandler.Stop()
		n.inprocHandler = nil
	}
}

// startIPC initializes and starts the IPC RPC endpoint.
func (n *Service) startIPC(apis []rpc.API) error {
	if n.ipcEndpoint == "" {
		return nil // IPC disabled.
	}
	listener, handler, err := rpc.StartIPCEndpoint(n.ipcEndpoint, apis)
	if err != nil {
		return err
	}
	n.ipcListener = listener
	n.ipcHandler = handler
	n.logger.Info("IPC endpoint opened", "url", n.ipcEndpoint)
	return nil
}

// stopIPC terminates the IPC RPC endpoint.
func (n *Service) stopIPC() {
	if n.ipcListener != nil {
		n.ipcListener.Close()
		n.ipcListener = nil

		n.logger.Info("IPC endpoint closed", "endpoint", n.ipcEndpoint)
	}
	if n.ipcHandler != nil {
		n.ipcHandler.Stop()
		n.ipcHandler = nil
	}
}

// startHTTP initializes and starts the HTTP RPC endpoint.
func (n *Service) startHTTP(endpoint string, apis []rpc.API, modules []string, cors []string, vhosts []string, timeouts rpc.HTTPTimeouts) error {
	// Short circuit if the HTTP endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts, timeouts)
	if err != nil {
		return err
	}
	n.logger.Info("HTTP endpoint opened", "url", fmt.Sprintf("http://%s", endpoint), "cors", strings.Join(cors, ","), "vhosts", strings.Join(vhosts, ","))
	// All listeners booted successfully
	n.httpEndpoint = endpoint
	n.httpListener = listener
	n.httpHandler = handler

	return nil
}

// stopHTTP terminates the HTTP RPC endpoint.
func (n *Service) stopHTTP() {
	if n.httpListener != nil {
		n.httpListener.Close()
		n.httpListener = nil

		n.logger.Info("HTTP endpoint closed", "url", fmt.Sprintf("http://%s", n.HTTPEndpoint()))
	}
	if n.httpHandler != nil {
		n.httpHandler.Stop()
		n.httpHandler = nil
	}
}

// startWS initializes and starts the websocket RPC endpoint.
func (n *Service) startWS(endpoint string, apis []rpc.API, modules []string, wsOrigins []string, exposeAll bool) error {
	// Short circuit if the WS endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartWSEndpoint(endpoint, apis, modules, wsOrigins, exposeAll)
	if err != nil {
		return err
	}
	n.logger.Info("WebSocket endpoint opened", "url", fmt.Sprintf("ws://%s", listener.Addr()))
	// All listeners booted successfully
	n.wsEndpoint = endpoint
	n.wsListener = listener
	n.wsHandler = handler

	return nil
}

// stopWS terminates the websocket RPC endpoint.
func (n *Service) stopWS() {
	if n.wsListener != nil {
		n.wsListener.Close()
		n.wsListener = nil

		n.logger.Info("WebSocket endpoint closed", "url", fmt.Sprintf("ws://%s", n.wsEndpoint))
	}
	if n.wsHandler != nil {
		n.wsHandler.Stop()
		n.wsHandler = nil
	}
}

// IPCEndpoint retrieves the current IPC endpoint used by the protocol stack.
func (n *Service) IPCEndpoint() string {
	return n.ipcEndpoint
}

// HTTPEndpoint retrieves the current HTTP endpoint used by the protocol stack.
func (n *Service) HTTPEndpoint() string {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.httpListener != nil {
		return n.httpListener.Addr().String()
	}
	return n.httpEndpoint
}

// WSEndpoint retrieves the current WS endpoint used by the protocol stack.
func (n *Service) WSEndpoint() string {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.wsListener != nil {
		return n.wsListener.Addr().String()
	}
	return n.wsEndpoint
}

// apis returns the collection of RPC descriptors this node offers.
func (n *Service) apis() []rpc.API {

	apis := GetNodeAPIs(n)
	apis = append(apis, GetEthAPIs(n)...)
	// TODO: turn on after implementation
	//var backend ethapi.Backend
	//apis = append(apis, ethapi.GetAPIs(backend)...)

	return apis
}
