package rpc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	delay = time.Second * 10
)
type RegistryDiscovery struct {
	*LocalDiscovery
	registryAddr   string
	delay    time.Duration
	updateAt time.Time
}


func NewRegistryDiscovery(registryAddr string, delay time.Duration)  *RegistryDiscovery{
	rd := &RegistryDiscovery{
		LocalDiscovery: NewLocalDiscovery(make([]string, 0)),
		registryAddr:      registryAddr,
		delay:       delay,
	}

	return rd
}

func (rd *RegistryDiscovery) Update(servers []string) error {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	rd.servers = servers
	rd.updateAt = time.Now()
	return nil
}

func (rd *RegistryDiscovery) Refresh() error {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	if rd.updateAt.Add(rd.delay).After(time.Now()) {
		return nil
	}

	log.Println("rpc registry: refresh servers from registry addr", rd.registryAddr)
	resp, err := http.Get(rd.registryAddr + "/pull")
	if err != nil {
		log.Println("rpc registry refresh err:", err)
		return err
	}

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		var nodes []string
		err = json.Unmarshal(body,&nodes)

		if err != nil {
			log.Println("rpc json decode response body err:", err)
			return err
		}

		rd.servers = make([]string, 0, len(nodes))
		rd.servers = append(rd.servers, nodes...)

		rd.updateAt = time.Now()

		return nil

	}else{
		return errors.New("rpc registry http.Get StatusCode not 200")
	}
}

func (rd *RegistryDiscovery) Get(loadmode LoadMode) (string, error) {
	if err := rd.Refresh(); err != nil {
		return "", err
	}

	return rd.LocalDiscovery.Get(loadmode)
}

func (rd *RegistryDiscovery) GetAll() ([]string, error) {

	if err := rd.Refresh(); err != nil {
		return nil, err
	}
	return rd.LocalDiscovery.GetAll()
}


