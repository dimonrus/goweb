package goweb

import (
	"github.com/dimonrus/gocli"
	"net"
	"net/http"
	"sync"
)

// Web config
type Config struct {
	Port    int
	Host    string
	Url     string
	Timeout struct {
		Read  int
		Write int
		Idle  int
	}
}

// Connection unique identifier
type ConnectionIdentifier string

// Event Name
type ConnectionEventName string

// Events that can be
type ConnectionEvent struct {
	Name ConnectionEventName
	Done chan bool
}

// Listeners struct
type ConnectionEventListeners struct {
	rw       sync.RWMutex
	listener map[ConnectionIdentifier][]ConnectionEvent
}

// Connections
type Connections struct {
	rw          sync.RWMutex
	connections map[ConnectionIdentifier]net.Conn
}

// Identifier for bind with connection
type BindingIdentifier interface{}

// Bindings
type ConnectionBindings struct {
	rw                  sync.RWMutex
	connectionBindingId map[ConnectionIdentifier]BindingIdentifier
	bindingIdConnection map[BindingIdentifier]ConnectionIdentifier
}

// Web application
type application struct {
	config Config
	app    gocli.Application
	server *http.Server
}
