package dependencyinjection

import (
	"fmt"
	"github.com/stretchr/testify/require"

	"net"
	"net/http"
	"testing"
)

func TestContainer(t *testing.T) {

	var HTTPBundle = Bundle(
		Provide(ProvideAddr("0.0.0.0", "8080")),
		Provide(NewMux, As(new(http.Handler))),
		Provide(NewHTTPServer, Prototype(), WithName("server")),
	)

	c := New(HTTPBundle)

	var server1 *http.Server
	err := c.Extract(&server1, Name("server"))
	require.NoError(t, err)

	var server2 *http.Server
	err = c.Extract(&server2, Name("server"))
	require.NoError(t, err)

	err = c.Invoke(PrintAddr)
	require.NoError(t, err)
}

// Addr
type Addr string

// ProvideAddr
func ProvideAddr(host string, port string) func() Addr {
	return func() Addr {
		return Addr(net.JoinHostPort(host, port))
	}
}

// NewHTTPServer
func NewHTTPServer(addr Addr, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    string(addr),
		Handler: handler,
	}
}

// NewMux
func NewMux() *http.ServeMux {
	return &http.ServeMux{}
}

// PrintAddr
func PrintAddr(addr Addr) {
	fmt.Println(addr)
}
