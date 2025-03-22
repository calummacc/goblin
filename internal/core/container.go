package core

import (
	"errors"
	"reflect"
	"sync"
)

var ErrNotFound = errors.New("dependency not found")

type Container struct {
	mutex      sync.RWMutex
	containers map[reflect.Type]interface{}
}

func NewContainer() *Container {	
	return &Container{
		containers: make(map[reflect.Type]interface{}),
	}
}

func (c *Container) Bind(interfaceType reflect.Type, implementation interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.containers[interfaceType] = implementation
}

func (c *Container) Resolve(interfaceType reflect.Type) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if implementation, exists := c.containers[interfaceType]; exists {
		return implementation, nil
	}
	return nil, ErrNotFound
}
