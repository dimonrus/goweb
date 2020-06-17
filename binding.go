package goweb

import "sync"

// Identifier for bind with connection
type BindingIdentifier interface{}

// Bindings
type ConnectionBindings struct {
	rw                  sync.RWMutex
	connectionBindingId map[ConnectionIdentifier]BindingIdentifier
	bindingIdConnection map[BindingIdentifier]ConnectionIdentifier
}

// Associate connection with id
func (cb *ConnectionBindings) Bind(bId BindingIdentifier, id ConnectionIdentifier) *ConnectionBindings {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	cb.bindingIdConnection[bId] = id
	cb.connectionBindingId[id] = bId
	return cb
}

// Unbind connection with id
func (cb *ConnectionBindings) UnBind(bId BindingIdentifier, id ConnectionIdentifier) *ConnectionBindings {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	delete(cb.bindingIdConnection, bId)
	delete(cb.connectionBindingId, id)
	return cb
}

// Get connection Id by id
func (cb *ConnectionBindings) GetConnectionId(bId BindingIdentifier) ConnectionIdentifier {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	return cb.bindingIdConnection[bId]
}

// Get participantId by connection id
func (cb *ConnectionBindings) GetBindingId(id ConnectionIdentifier) BindingIdentifier {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	return cb.connectionBindingId[id]
}

// Get all bindings identifiers
func (cb *ConnectionBindings) GetBindingIdentifiers() []BindingIdentifier {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	result := make([]BindingIdentifier, 0)
	for id, _ := range cb.bindingIdConnection {
		result = append(result, id)
	}
	return result
}

// New bindings
func NewConnectionBindings() *ConnectionBindings {
	return &ConnectionBindings{
		connectionBindingId: make(map[ConnectionIdentifier]BindingIdentifier),
		bindingIdConnection: make(map[BindingIdentifier]ConnectionIdentifier),
	}
}
