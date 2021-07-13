package dependencyinjection

import (
	"github.com/magiconair/properties/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type A struct {
	Name string
}

var ls string

func NewA() *A {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ls = "A-" + strconv.Itoa(r.Int())
	return &A{Name: ls}
}

func Test_DI_Register(t *testing.T) {
	services := NewServiceCollection()

	services.AddSingleton(NewA)

	serviceProvider := services.Build()

	var env *A
	_ = serviceProvider.GetService(&env)
	assert.Equal(t, env.Name, ls)

	_ = serviceProvider.GetService(&env)
	assert.Equal(t, env.Name, ls)

	_ = serviceProvider.GetService(&env)
	assert.Equal(t, env.Name, ls)
}
