package goweb

import (
	"net"
	"sync"
)

// Connection unique identifier
type ConnectionIdentifier string

// Connections
type Connections struct {
	rw          sync.RWMutex
	connections map[ConnectionIdentifier]net.Conn
}

// Get connection if exists
func (cs *Connections) Get(id ConnectionIdentifier) net.Conn {
	cs.rw.RLock()
	defer cs.rw.RUnlock()
	return cs.connections[id]
}

// Set connection
func (cs *Connections) Set(id ConnectionIdentifier, conn net.Conn) *Connections {
	cs.rw.Lock()
	defer cs.rw.Unlock()
	cs.connections[id] = conn
	return cs
}

// Unset connection if exists
func (cs *Connections) Unset(id ConnectionIdentifier) *Connections {
	cs.rw.Lock()
	defer cs.rw.Unlock()
	if _, ok := cs.connections[id]; ok {
		delete(cs.connections, id)
	}
	return cs
}

// Connections len
func (cs *Connections) Len() int {
	cs.rw.Lock()
	defer cs.rw.Unlock()
	return len(cs.connections)
}

// Init tcp connection pool
func NewConnections() *Connections {
	return &Connections{
		connections: make(map[ConnectionIdentifier]net.Conn),
	}
}
