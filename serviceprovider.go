package dependencyinjection

type IServiceProvider interface {
	GetService(refObject interface{}) error
	GetServiceByName(refObject interface{}, name string) error
	GetServiceByTags(refObject interface{}, tags map[string]string) (err error)
	GetGraph() string
	InvokeService(fn interface{}) error
}
