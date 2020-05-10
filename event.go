package goweb

import "sync"

// Event Name
type ConnectionEventName string

// Events that can be
type ConnectionEvent struct {
	// Name of the event
	Name ConnectionEventName
	// Mark when event happens
	Done chan bool
}

// Listeners struct
type ConnectionEventListeners struct {
	rw       sync.RWMutex
	listener map[ConnectionIdentifier][]ConnectionEvent
}

// Register connection listener
func (cel *ConnectionEventListeners) Register(id ConnectionIdentifier, name ConnectionEventName) <-chan bool {
	cel.rw.Lock()
	defer cel.rw.Unlock()
	event := &ConnectionEvent{Name: name, Done: make(chan bool, 1)}
	cel.listener[id] = append(cel.listener[id], *event)
	return event.Done
}

// Unregister connection event listener
func (cel *ConnectionEventListeners) Unregister(id ConnectionIdentifier, name ConnectionEventName) *ConnectionEventListeners {
	cel.rw.Lock()
	defer cel.rw.Unlock()
	for i, value := range cel.listener[id] {
		if value.Name == name {
			cel.listener[id] = append(cel.listener[id][:i], cel.listener[id][i+1:]...)
		}
	}
	return cel
}

// Unregister connection listener
func (cel *ConnectionEventListeners) UnregisterConnection(id ConnectionIdentifier) *ConnectionEventListeners {
	cel.rw.Lock()
	defer cel.rw.Unlock()
	delete(cel.listener, id)
	return cel
}

// Dispatch event
func (cel *ConnectionEventListeners) Dispatch(id ConnectionIdentifier, name ConnectionEventName) *ConnectionEventListeners {
	cel.rw.RLock()
	defer cel.rw.RUnlock()
	if events, ok := cel.listener[id]; ok {
		for i := range events {
			if events[i].Name == name && len(events[i].Done) == 0 {
				events[i].Done <- true
			}
		}
	}
	return cel
}

// Get events
func (cel *ConnectionEventListeners) Get(id ConnectionIdentifier) []ConnectionEvent {
	cel.rw.RLock()
	defer cel.rw.RUnlock()
	return cel.listener[id]
}

// New Connection event Listener
func NewConnectionEventListeners() *ConnectionEventListeners {
	return &ConnectionEventListeners{listener: make(map[ConnectionIdentifier][]ConnectionEvent)}
}
