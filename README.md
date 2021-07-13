# Dependency injection
Dependency injection for Go programming language.

Dependency injection is one form of the broader technique of inversion of control. It is used to increase modularity of the program and make it extensible.

## Examples
```go
type A struct {
	Name string
}

func NewA() *A {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	name := "A-" + strconv.Itoa(r.Int())
	return &A{Name: ls}
}

services := NewServiceCollection()
services.AddSingleton(NewA)
//serviceCollection.AddSingletonByImplementsAndName("redis-master", NewRedis, new(abstractions.IDataSource))
//serviceCollection.AddTransientByImplements(NewRedisClient, new(redis.IClient))
//serviceCollection.AddTransientByImplements(NewRedisHealthIndicator, new(health.Indicator))
serviceProvider := services.Build()

var env *A
_ = serviceProvider.GetService(&env) // used
```

## How will dependency injection help me?

Dependency injection is one form of the broader technique of inversion
of control. It is used to increase modularity of the program and make it
extensible.

## Contents

- [Installing](#installing)
- [Tutorial](#tutorial)
    - [Providing](#providing)
    - [Extraction](#extraction)
    - [Invocation](#invocation)
    - [Lazy-loading](#lazy-loading)
    - [Interfaces](#interfaces)
    - [Groups](#groups)
- [Advanced features](#advanced-features)
    - [Named definitions](#named-definitions)
    - [Optional parameters](#optional-parameters)
    - [Parameter Bag](#parameter-bag)
    - [Prototypes](#prototypes)
    - [Cleanup](#cleanup)
    - [Visualization](#visualization)
- [Contributing](#contributing)

## Installing

```shell
go get -u github.com/yoyofxteam/dependencyinjection@v1.0.0
```

### Providing

To start, we will need to create two fundamental types: `http.Server`
and `http.ServeMux`. Let's create a simple constructors that initialize
it:

```go
// NewServer creates a http server with provided mux as handler.
func NewServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Handler: mux,
	}
}

// NewServeMux creates a new http serve mux.
func NewServeMux() *http.ServeMux {
	return &http.ServeMux{}
}
```

> Supported constructor signature:
>
> ```go
> func([dep1, dep2, depN]) (result, [cleanup, error])
> ```

Now let's teach a container to build these types.

```go
import (
  di "github.com/yoyofxteam/dependencyinjection"
)

container := di.New(
	// provide http server
    di.Provide(NewServer),
    // provide http serve mux
    di.Provide(NewServeMux)
)
```

The function `di.New()` parse our constructors, compile dependency
graph and return `*di.Container` type for interaction. Container
panics if it could not compile.

> I think that panic at the initialization of the application and not in
> runtime is usual.

### Extraction

We can extract the built server from the container. For this, define the
variable of extracted type and pass variable pointer to `Extract`
function.

> If extracted type not found or the process of building instance cause
> error, `Extract` return error.

If no error occurred, we can use the variable as if we had built it
yourself.

```go
// declare type variable
var server *http.Server
// extracting
err := container.Extract(&server)
if err != nil {
	// check extraction error
}

server.ListenAndServe()
```

> Note that by default, the container creates instances as a singleton.
> But you can change this behaviour. See [Prototypes](#prototypes).

### Invocation

As an alternative to extraction we can use `Invoke()` function. It
resolves function dependencies and call the function. Invoke function
may return optional error.

```go
// StartServer starts the server.
func StartServer(server *http.Server) error {
    return server.ListenAndServe()
}

container.Invoke(StartServer)
```

### Lazy-loading

Result dependencies will be lazy-loaded. If no one requires a type from
the container it will not be constructed.

### Interfaces

Inject make possible to provide implementation as an interface.

```go
// NewServer creates a http server with provided mux as handler.
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
	}
}
```

For a container to know that as an implementation of `http.Handler` is
necessary to use, we use the option `di.As()`. The arguments of this
option must be a pointer(s) to an interface like `new(Endpoint)`.

> This syntax may seem strange, but I have not found a better way to
> specify the interface.

Updated container initialization code:

```go
container := inject.New(
	// provide http server
	inject.Provide(NewServer),
	// provide http serve mux as http.Handler interface
	inject.Provide(NewServeMux, inject.As(new(http.Handler)))
)
```

Now container uses provide `*http.ServeMux` as `http.Handler` in server
constructor. Using interfaces contributes to writing more testable code.

### Groups

Container automatically groups all implementations of interface to
`[]<interface>` group. For example, provide with
`inject.As(new(http.Handler)` automatically creates a group
`[]http.Handler`.

Let's add some http controllers using this feature. Controllers have
typical behavior. It is registering routes. At first, will create an
interface for it.

```go
// Controller is an interface that can register its routes.
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}
```

Now we will write controllers and implement `Controller` interface.

##### OrderController

```go
// OrderController is a http controller for orders.
type OrderController struct {}

// NewOrderController creates a auth http controller.
func NewOrderController() *OrderController {
	return &OrderController{}
}

// RegisterRoutes is a Controller interface implementation.
func (a *OrderController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/orders", a.RetrieveOrders)
}

// Retrieve loads orders and writes it to the writer.
func (a *OrderController) RetrieveOrders(writer http.ResponseWriter, request *http.Request) {
	// implementation
}
```

##### UserController

```go
// UserController is a http endpoint for a user.
type UserController struct {}

// NewUserController creates a user http endpoint.
func NewUserController() *UserController {
	return &UserController{}
}

// RegisterRoutes is a Controller interface implementation.
func (e *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", e.RetrieveUsers)
}

// Retrieve loads users and writes it using the writer.
func (e *UserController) RetrieveUsers(writer http.ResponseWriter, request *http.Request) {
    // implementation
}
```

Just like in the example with interfaces, we will use `inject.As()`
provide option.

```go
container := inject.New(
	di.Provide(NewServer),        // provide http server
	di.Provide(NewServeMux),       // provide http serve mux
	// endpoints
	di.Provide(NewOrderController, di.As(new(Controller))),  // provide order controller
	di.Provide(NewUserController, di.As(new(Controller))),  // provide user controller
)
```

Now, we can use `[]Controller` group in our mux. See updated code:

```go
// NewServeMux creates a new http serve mux.
func NewServeMux(controllers []Controller) *http.ServeMux {
	mux := &http.ServeMux{}

	for _, controller := range controllers {
		controller.RegisterRoutes(mux)
	}

	return mux
}
```

## Advanced features

### Named definitions

In some cases you have more than one instance of one type. For example
two instances of database: master - for writing, slave - for reading.

First way is a wrapping types:

```go
// MasterDatabase provide write database access.
type MasterDatabase struct {
	*Database
}

// SlaveDatabase provide read database access.
type SlaveDatabase struct {
	*Database
}
```

Second way is a using named definitions with `di.WithName()` provide
option:

```go
// provide master database
di.Provide(NewMasterDatabase, di.WithName("master"))
// provide slave database
di.Provide(NewSlaveDatabase, di.WithName("slave"))
```

If you need to extract it from container use `di.Name()` extract
option.

```go
var db *Database
container.Extract(&db, di.Name("master"))
```

If you need to provide named definition in other constructor use
`di.Parameter` with embedding.

```go
// ServiceParameters
type ServiceParameters struct {
	di.Parameter
	
	// use `di` tag for the container to know that field need to be injected.
	MasterDatabase *Database `di:"master"`
	SlaveDatabase *Database  `di:"slave"`
}

// NewService creates new service with provided parameters.
func NewService(parameters ServiceParameters) *Service {
	return &Service{
		MasterDatabase:  parameters.MasterDatabase,
		SlaveDatabase: parameters.SlaveDatabase,
	}
}
```

### Optional parameters

Also `di.Parameter` provide ability to skip dependency if it not exists
in container.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Parameter
	
	Logger *Logger `di:"optional"`
}
```

> Constructors that declare dependencies as optional must handle the
> case of those dependencies being absent.

You can use naming and optional together.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Parameter
	
	StdOutLogger *Logger `di:"stdout"`
	FileLogger   *Logger `di:"file,optional"`
}
```

### Parameter Bag

If you need to specify some parameters on definition level you can use
`inject.ParameterBag` provide option. This is a `map[string]interface{}`
that transforms to `di.ParameterBag` type.

```go
// Provide server with parameter bag
di.Provide(NewServer, di.ParameterBag{
	"addr": ":8080",
})

// NewServer create a server with provided parameter bag. Note: use di.ParameterBag type.
// Not inject.ParameterBag.
func NewServer(pb di.ParameterBag) *http.Server {
	return &http.Server{
		Addr: pb.RequireString("addr"),
	}
}
```

### Prototypes

If you want to create a new instance on each extraction use
`di.Prototype()` provide option.

```go
di.Provide(NewRequestContext, di.Prototype())
```

> todo: real use case

### Cleanup

If a provider creates a value that needs to be cleaned up, then it can
return a closure to clean up the resource.

```go
func NewFile(log Logger, path Path) (*os.File, func(), error) {
    f, err := os.Open(string(path))
    if err != nil {
        return nil, nil, err
    }
    cleanup := func() {
        if err := f.Close(); err != nil {
            log.Log(err)
        }
    }
    return f, cleanup, nil
}
```

After `container.Cleanup()` call, it iterate over instances and call
cleanup function if it exists.

```go
container := di.New(
	// ...
    di.Provide(NewFile),
)

// do something
container.Cleanup() // file was closed
```

> Cleanup now work incorrectly with prototype providers.

## Visualization

Dependency graph may be presented via
([Graphviz](https://www.graphviz.org/)). For it, load string
representation:

```go
var graph *di.di.Graph
if err = container.Extract(&graph); err != nil {
    // handle err
}

dotGraph := graph.String() // use string representation
```

And paste it to <a href="https://dreampuf.github.io/GraphvizOnline"
target="_blank">graphviz online tool</a>:

<img src="https://github.com/defval/inject/raw/master/graph.png">