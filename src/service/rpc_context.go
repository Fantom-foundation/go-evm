package service

import (
	"reflect"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

// RpcServiceContext is a collection of service independent options inherited from
// the protocol stack, that is passed to all constructors to be optionally used;
// as well as utility methods to operate on the service environment.
type RpcServiceContext struct {
	config         *node.Config
	services       map[reflect.Type]RpcService // Index of the already constructed services
	EventMux       *event.TypeMux              // Event multiplexer used for decoupled notifications
	AccountManager *accounts.Manager           // Account manager created by the node.
}

// RpcService retrieves a currently running service registered of a specific type.
func (ctx *RpcServiceContext) Service(service interface{}) error {
	element := reflect.ValueOf(service).Elem()
	if running, ok := ctx.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return ErrServiceUnknown
}

// RpcServiceConstructor is the function signature of the constructors needed to be
// registered for service instantiation.
type RpcServiceConstructor func(ctx *RpcServiceContext) (RpcService, error)

// Service is an individual protocol that can be registered into a node.
//
// Notes:
//
// • Service life-cycle management is delegated to the node. The service is allowed to
// initialize itself upon creation, but no goroutines should be spun up outside of the
// Start method.
//
// • Restart logic is not required as the node will create a fresh instance
// every time a service is started.
type RpcService interface {

	// APIs retrieves the list of RPC descriptors the service provides
	APIs() []rpc.API

	// Start is called after all services have been constructed and the networking
	// layer was also initialized to spawn any goroutines required by the service.
	Start() error

	// Stop terminates all goroutines belonging to the service, blocking until they
	// are all terminated.
	Stop() error
}
