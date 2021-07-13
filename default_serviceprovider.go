package dependencyinjection

import (
	"github.com/yoyofxteam/dependencyinjection/di"
)

type DefaultServiceProvider struct {
	container *Container
}

func (d DefaultServiceProvider) GetService(refObject interface{}) (err error) {
	err = d.container.Extract(refObject)
	return err
}

func (d DefaultServiceProvider) GetServiceByName(refObject interface{}, name string) (err error) {
	err = d.container.Extract(refObject, Name(name))

	return err
}

func (d DefaultServiceProvider) GetGraph() string {
	var graph *di.Graph
	if err := d.container.Extract(&graph); err != nil {
		// handle err
	}

	return graph.String() // use string representation
}

func (d DefaultServiceProvider) InvokeService(fn interface{}) error {
	return d.container.Invoke(fn)
}
