package dependencyinjection

import (
	"github.com/goava/di"
	"unsafe"
)

type DefaultServiceProvider struct {
	container *di.Container
}

func (d DefaultServiceProvider) GetService(refObject interface{}) (err error) {
	err = d.container.Resolve(refObject)
	return err
}

func (d DefaultServiceProvider) GetServiceByName(refObject interface{}, name string) (err error) {
	err = d.container.Resolve(refObject, di.Name(name))
	return err
}

func (d DefaultServiceProvider) GetServiceByTags(refObject interface{}, tags map[string]string) (err error) {
	p := unsafe.Pointer(&tags)
	var tag di.Tags
	tag = *(*di.Tags)(p)
	err = d.container.Resolve(refObject, tag)
	return err
}

func (d DefaultServiceProvider) GetGraph() string {
	return ""
}

func (d DefaultServiceProvider) InvokeService(fn interface{}) error {
	return d.container.Invoke(fn)
}
