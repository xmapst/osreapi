package dag

import (
	"fmt"
	"sync"
)

var manager = sync.Map{}

const (
	graphPrefix  = "graph#%s#graph"
	vertexPrefix = "graph#%s#vertex#%s#vertex#graph"
)

type IManager interface {
	Name() string
	Kill() error
	Pause(duration string) error
	Resume()
	Paused() bool
}

func leave(key string) (IManager, error) {
	value, ok := manager.Load(key)
	if !ok {
		return nil, ErrNotFound
	}
	m, ok := value.(IManager)
	if !ok {
		manager.Delete(key)
		return nil, ErrWrongType
	}
	return m, nil
}

func join(key string, iManager IManager) {
	manager.Store(key, iManager)
}

func remove(key string) {
	manager.Delete(key)
}

func GraphManager(gName string) (IManager, error) {
	return leave(fmt.Sprintf(graphPrefix, gName))
}

func VertexManager(gName, vName string) (IManager, error) {
	return leave(fmt.Sprintf(vertexPrefix, gName, vName))
}
