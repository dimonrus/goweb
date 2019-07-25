package goweb

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
	if connId, ok := cb.bindingIdConnection[bId]; ok {
		return connId
	}

	return ""
}

// Get participantId by connection id
func (cb *ConnectionBindings) GetBindingId(id ConnectionIdentifier) BindingIdentifier {
	cb.rw.Lock()
	defer cb.rw.Unlock()
	if id, ok := cb.connectionBindingId[id]; ok {
		return id
	}

	return nil
}

// New bindings
func NewConnectionBindings() *ConnectionBindings {
	return &ConnectionBindings{
		connectionBindingId: make(map[ConnectionIdentifier]BindingIdentifier),
		bindingIdConnection: make(map[BindingIdentifier]ConnectionIdentifier),
	}
}
