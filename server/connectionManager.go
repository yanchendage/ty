package server

import (
	"errors"
	"sync"
)

type ConnectionManager struct {
	connections map[uint32]IConnection
	mu sync.RWMutex
}

func NewConnectionManager() *ConnectionManager{
	return  &ConnectionManager{
		connections: make(map[uint32]IConnection),
	}
}

func (cm *ConnectionManager) Add(connection IConnection)  {
	cm.mu.Lock()
	defer  cm.mu.Unlock()

	cm.connections[connection.GetConnID()] = connection
}

func (cm *ConnectionManager) Remove(connectionID uint32)  {
	cm.mu.Lock()
	defer  cm.mu.Unlock()

	delete(cm.connections, connectionID)
}

func (cm *ConnectionManager) Get(connectionID uint32) (IConnection, error) {
	cm.mu.RLock()
	defer  cm.mu.RUnlock()

	connection, ok := cm.connections[connectionID]

	if !ok {
		return nil, errors.New("【Connection Manager】connection not found")
	}else {
		return connection, nil
	}
}

func (cm *ConnectionManager) Count() int {
	cm.mu.Lock()
	defer  cm.mu.Unlock()

	return len(cm.connections)
}

func (cm *ConnectionManager) Clean() {
	cm.mu.Lock()
	defer  cm.mu.Unlock()

	for ID, connection := range cm.connections{
		delete(cm.connections, ID)
		connection.Stop()
	}
}





