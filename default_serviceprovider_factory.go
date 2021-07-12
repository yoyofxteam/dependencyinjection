package dependencyinjection

import (
	"github.com/goava/di"
)

func (sc ServiceCollection) Build() IServiceProvider {
	var providers []di.Option
	for _, desc := range sc.serviceDescriptors {

		var providerOptions []di.ProvideOption
		if desc.Implements != nil {
			providerOptions = append(providerOptions, di.As(desc.Implements))
		}
		if desc.Name != "" {
			providerOptions = append(providerOptions, di.WithName(desc.Name))
		}
		//prototype is create a new instance on each call.
		if desc.Lifetime != Singleton {
			//providerOptions = append(providerOptions)
		}

		provider := di.Provide(desc.Provider, providerOptions...)
		providers = append(providers, provider)

	}
	container, _ := di.New(providers...)

	return &DefaultServiceProvider{container}
}
