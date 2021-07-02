package rpc

import (
	"testing"
)



func TestServer(t *testing.T) {
	var foo Foo
	NewServer("RPC", "127.0.0.1", 7727)

	sm := NewServiceManager()
	sm.RegisterService(&foo)
}


func TestNewServerManager(t *testing.T) {
	var foo Foo
	serviceManager := NewServiceManager()
	serviceManager.RegisterService(&foo)

	properties := map[string]interface{}{
		"ServiceManager" : serviceManager,
	}

	NewServerManager("RPC", "127.0.0.1", 7729, properties)
}

func TestInitServerManager(t *testing.T) {
	var foo Foo

	serverManager := InitServerManager("RPC", "127.0.0.1", 7729,"http://127.0.0.1:8888")
	serverManager.RegisterService(&foo)

	serverManager.Run()
}