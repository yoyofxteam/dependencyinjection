package dependencyinjection

import (
	"github.com/goava/di"
	"github.com/magiconair/properties/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"unsafe"
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

func Test_MapConvertTags(t *testing.T) {
	tags := map[string]string{"name": "A1"}
	p := unsafe.Pointer(&tags)
	var tag di.Tags
	tag = *(*di.Tags)(p)
	assert.Equal(t, tag != nil, true)
}
