package docker

import (
	"fmt"
	"sync"

	"github.com/docker/docker/api/types"
)

type Containers struct {
	byListen map[string]*types.ContainerJSON
	byId     map[string]*types.ContainerJSON
	lock     *sync.RWMutex
}

func NewContainers() *Containers {
	return &Containers{
		byListen: make(map[string]*types.ContainerJSON),
		byId:     make(map[string]*types.ContainerJSON),
		lock:     &sync.RWMutex{},
	}
}

func (c *Containers) Add(container *types.ContainerJSON) bool {
	if len(container.NetworkSettings.Networks) == 0 { // the container has no IP
		return false
	}
	if len(container.NetworkSettings.Ports) == 0 { // the container expose nothing
		return false
	}
	ips := make([]string, 0)
	for _, net := range container.NetworkSettings.Networks {
		ips = append(ips, net.IPAddress)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	for port := range container.NetworkSettings.Ports {
		for _, ip := range ips {
			// FIXME handling port.Rangge
			c.byListen[fmt.Sprintf("%s:%s", ip, port.Port())] = container
		}
	}
	c.byId[container.ID] = container
	return true
}

func (c *Containers) Remove(id string) bool {
	c.lock.RLock()
	_, ok := c.byId[id]
	c.lock.RUnlock()
	if !ok {
		return false
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.byId, id)
	tombstones := make([]string, 0)
	for listen, container := range c.byListen {
		if container.ID == id {
			tombstones = append(tombstones, listen)
		}
	}
	for _, tombstone := range tombstones {
		delete(c.byListen, tombstone)
	}
	return true
}

func (c *Containers) ByListen(listen string) *types.ContainerJSON {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.byListen[listen]
}
