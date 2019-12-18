// Modified from go-ethereum under GNU Lesser General Public License
package service

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/QuarkChain/goquarkchain/p2p"
	"github.com/QuarkChain/goquarkchain/rpc"
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

// ServiceContext is a collection of service independent options inherited from
// the protocol stack, that is passed to all constructors to be optionally used;
// as well as utility methods to operate on the service environment.
type ServiceContext struct {
	config    *Config
	services  map[reflect.Type]Service // Index of the already constructed services
	Shutdown  chan os.Signal
	Timestamp time.Time
	EventMux  *event.TypeMux // Event multiplexer used for decoupled notifications
}

func (ctx *ServiceContext) WSIsAlive() bool {
	return ctx.config.WSEndpoint != ""
}

// makeDatabaseHandles raises out the number of allowed file handles per process
// for Geth and returns half of the allowance to assign to the database.
func makeDatabaseHandles() int {
	limit, err := fdlimit.Maximum()
	if err != nil {
		panic(fmt.Errorf("Failed to retrieve file descriptor allowance: %v", err))
	}
	if err := fdlimit.Raise(uint64(limit)); err != nil {
		panic(fmt.Errorf("Failed to raise file descriptor allowance: %v", err))
	}
	return limit / 2 // Leave half for networking and other stuff
}

// OpenDatabase opens an existing database with the given name (or creates one
// if no previous can be found) from within the node's data directory. If the
// node is an ephemeral one, a memory database is returned.
func (ctx *ServiceContext) OpenDatabase(name string, clean bool, isReadOnly bool) (ethdb.Database, error) {
	if ctx.config == nil || ctx.config.DataDir == "" {
		return NewQkcMemoryDB(isReadOnly), nil
	}
	//db, err := qkcdb.NewRDBDatabase(ctx.config.ResolvePath(name), clean, isReadOnly)
	//if err != nil {
	//	return nil, err
	//}
	db, err := ethdb.NewLDBDatabase(name, 32, 64)
	return db, err
}

// ResolvePath resolves a user path into the data directory if that was relative
// and if the user actually uses persistent storage. It will return an empty string
// for emphemeral storage and the user's own input for absolute paths.
func (ctx *ServiceContext) ResolvePath(path string) string {
	return ctx.config.ResolvePath(path)
}

// Service retrieves a currently running service registered of a specific type.
func (ctx *ServiceContext) Service(service interface{}) error {
	element := reflect.ValueOf(service).Elem()
	if running, ok := ctx.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return ErrServiceUnknown
}

// ServiceConstructor is the function signature of the constructors needed to be
// registered for service instantiation.
type ServiceConstructor func(ctx *ServiceContext) (Service, error)

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
type Service interface {
	// Protocols retrieves the P2P protocols the service wishes to start.
	Protocols() []p2p.Protocol

	// APIs retrieves the list of RPC descriptors the service provides
	APIs() []rpc.API

	// Start is called after all services have been constructed and the networking
	// layer was also initialized to spawn any goroutines required by the service.
	Init(server *p2p.Server) error

	// Stop terminates all goroutines belonging to the service, blocking until they
	// are all terminated.
	Stop() error
}
