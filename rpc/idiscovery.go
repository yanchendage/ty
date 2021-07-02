package rpc

type LoadMode int

const (
	Random LoadMode = iota
	RoundRobin
)

type IDiscovery interface {
	Refresh() error
	Update(servers []string) error
	Get(mode LoadMode) (string, error)
	GetAll() ([]string,error)
}


