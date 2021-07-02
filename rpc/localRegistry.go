package rpc

import (
	"github.com/yanchendage/ty/web"
	"log"
	"sort"
	"sync"
	"time"
)

var lr *LocalRegistry

type LocalRegistry struct {
	timeout time.Duration
	mu sync.Mutex
	nodes map[string]*Node
}

type Node struct {
	addr string
	start time.Time
}

const (
	path = "/ty/rpc/registry"
	timeout = time.Minute * 5
)

func NewLocalRegistry(timeout time.Duration)  *LocalRegistry{
	return &LocalRegistry{
		timeout:  timeout,
		nodes: make(map[string]*Node),
	}
}

func (lr *LocalRegistry) push(addr string)  {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	node := lr.nodes[addr]

	if node == nil {
		lr.nodes[addr] = &Node{addr, time.Now()}
	}else {
		node.start = time.Now()
	}
}

func (lr *LocalRegistry) pull() []string {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	var addrs []string

	for addr, node := range lr.nodes {
		if lr.timeout == 0 || node.start.Add(lr.timeout).After(time.Now()) {
			addrs = append(addrs, addr)
		}else{
			delete(lr.nodes, addr)
		}
	}

	sort.Strings(addrs)
	return addrs
}

//var DefaultLocalRegistry = NewLocalRegistry(timeout)

func routerPull(c *web.Context)  {
	addrs := lr.pull()
	c.Json(200, addrs)
}

func routerPush(c *web.Context)  {
	if c.Query("addr") != "" {
		lr.push(c.Query("addr"))
	}
}

func NewLocalRegistryServer(addr string){
	lr = NewLocalRegistry(timeout)

	r := web.New()
	r.GET("/pull",routerPull)
	r.POST("/push",routerPush)

	r.Run(addr)

	log.Println("local registry server starting")
}




