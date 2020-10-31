package registry

import (
	"github.com/micro/micro/v3/service/registry"
)

//SortedServices for registry.Service
type SortedServices struct {
	services []*registry.Service
}

func (s SortedServices) Len() int {
	return len(s.services)
}

func (s SortedServices) Less(i, j int) bool {
	return s.services[i].Name < s.services[j].Name
}

func (s SortedServices) Swap(i, j int) {
	s.services[i], s.services[j] = s.services[j], s.services[i]
}
