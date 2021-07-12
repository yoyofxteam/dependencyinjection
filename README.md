# DI
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
serviceProvider := services.Build()

var env *A
_ = serviceProvider.GetService(&env) // used
```