package registry

import (
	"github.com/micro/micro/v3/service/registry"
)

//SortedServices for registry.Service
type SortedServices struct {
	Services []*registry.Service
}

func (s SortedServices) Len() int {
	return len(s.Services)
}

func (s SortedServices) Less(i, j int) bool {
	return s.Services[i].Name < s.Services[j].Name
}

func (s SortedServices) Swap(i, j int) {
	s.Services[i], s.Services[j] = s.Services[j], s.Services[i]
}
