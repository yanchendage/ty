package rpc

import (
	"context"
	"io"
	"sync"
)

type ClientManager struct {
	discovery IDiscovery
	loadMode LoadMode
	mu sync.Mutex
	clients map[string]*Client
}

var _ io.Closer = (*ClientManager)(nil)

func NewClientManager(discover IDiscovery, loadMode LoadMode) *ClientManager {
	return &ClientManager{
		discovery: discover,
		loadMode:  loadMode,
		clients: map[string]*Client{},
	}
}

func (cm *ClientManager) Close() error  {
	cm.mu.Lock()
	defer  cm.mu.Unlock()

	for key, client := range  cm.clients{
		client.Close()
		delete(cm.clients,key)
	}
	return nil
}

func (cm *ClientManager) dial(addr string) (*Client, error){
	cm.mu.Lock()
	defer  cm.mu.Unlock()

	client, ok := cm.clients[addr]
	if ok && !client.Available() {
		client.Close()
		delete(cm.clients, addr)
		client = nil
	}

	if client == nil {
		var err error
		//client, err = Dial(addr,"application/gob")
		client, err = Dial(addr,"application/protobuf")
		if err != nil {
			return nil, err
		}
		cm.clients[addr] = client
	}

	return client, nil
}

func (cm *ClientManager) call(addr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	client, err := cm.dial(addr)

	if err != nil {
		return err
	}

	return client.SyncCall(ctx, serviceMethod, args, reply)
}

func (cm *ClientManager) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {

	addr, err := cm.discovery.Get(cm.loadMode)

	if err != nil {
		return err
	}
	return cm.call(addr, ctx, serviceMethod, args, reply)
}


