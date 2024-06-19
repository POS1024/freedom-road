package inter

import "sync"

var ConfiguratorLocker sync.Mutex

type Configurator interface {
	Store() error
	Read() error
	Get(field string) (interface{}, bool)
	Set(field string, value interface{}) bool
}
