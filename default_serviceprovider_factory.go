package dependencyinjection

func (sc ServiceCollection) Build() IServiceProvider {
	var providers []Option
	for _, desc := range sc.serviceDescriptors {

		var providerOptions []ProvideOption
		if desc.Implements != nil {
			providerOptions = append(providerOptions, As(desc.Implements))
		}
		if desc.Name != "" {
			providerOptions = append(providerOptions, WithName(desc.Name))
		}
		//prototype is create a new instance on each call.
		if desc.Lifetime != Singleton {
			providerOptions = append(providerOptions, Prototype())
		}

		provider := Provide(desc.Provider, providerOptions...)
		providers = append(providers, provider)

	}
	container := New(providers...)

	return &DefaultServiceProvider{container}
}
