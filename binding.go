package goweb

import "sync"

// BindingIdentifier Identifier for bind with connection
type BindingIdentifier interface{}

// ConnectionBindings Bindings
type ConnectionBindings struct {
	rw                  sync.RWMutex
	connectionBindingId map[ConnectionIdentifier]BindingIdentifier
	bindingIdConnection map[BindingIdentifier]ConnectionIdentifier
}

// Bind Associate connection with id
func (cb *ConnectionBindings) Bind(bId BindingIdentifier, id ConnectionIdentifier) *ConnectionBindings {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	cb.bindingIdConnection[bId] = id
	cb.connectionBindingId[id] = bId
	return cb
}

// UnBind connection with id
func (cb *ConnectionBindings) UnBind(bId BindingIdentifier, id ConnectionIdentifier) *ConnectionBindings {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	delete(cb.bindingIdConnection, bId)
	delete(cb.connectionBindingId, id)
	return cb
}

// GetConnectionId Get connection identifier by id
func (cb *ConnectionBindings) GetConnectionId(bId BindingIdentifier) ConnectionIdentifier {
	cb.rw.RLock()
	defer cb.rw.RUnlock()
	return cb.bindingIdConnection[bId]
}

// GetBindingId Get participantId by connection id
func (cb *ConnectionBindings) GetBindingId(id ConnectionIdentifier) BindingIdentifier {
	cb.rw.RLock()
	defer cb.rw.RUnlock()
	return cb.connectionBindingId[id]
}

// GetBindingIdentifiers Get all bindings identifiers
func (cb *ConnectionBindings) GetBindingIdentifiers() []BindingIdentifier {
	cb.rw.RLock()
	defer cb.rw.RUnlock()
	result := make([]BindingIdentifier, 0)
	for id, _ := range cb.bindingIdConnection {
		result = append(result, id)
	}
	return result
}

// NewConnectionBindings New bindings
func NewConnectionBindings() *ConnectionBindings {
	return &ConnectionBindings{
		connectionBindingId: make(map[ConnectionIdentifier]BindingIdentifier),
		bindingIdConnection: make(map[BindingIdentifier]ConnectionIdentifier),
	}
}
