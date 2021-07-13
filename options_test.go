package dependencyinjection

import (
	"github.com/yoyofxteam/dependencyinjection/di"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestProvideOptions
func TestProvideOptions(t *testing.T) {
	opts := &di.ProvideParams{
		Parameters: map[string]interface{}{},
	}

	for _, opt := range []ProvideOption{
		WithName("test"),
		As(new(http.Handler)),
		Prototype(),
		ParameterBag{
			"test": "test",
		},
	} {
		opt.apply(opts)
	}

	require.Equal(t, &di.ProvideParams{
		Name:        "test",
		Interfaces:  []interface{}{new(http.Handler)},
		IsPrototype: true,
		Parameters: map[string]interface{}{
			"test": "test",
		},
	}, opts)
}

func TestExtractOptions(t *testing.T) {
	opts := &di.ExtractParams{}

	for _, opt := range []ExtractOption{
		Name("test"),
	} {
		opt.apply(opts)
	}

	require.Equal(t, &di.ExtractParams{
		Name: "test",
	}, opts)
}
