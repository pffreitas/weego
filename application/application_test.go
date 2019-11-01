package application_test

import (
	"testing"

	"github.com/pffreitas/weego/application"
	whttp "github.com/pffreitas/weego/application/http"
	"github.com/stretchr/testify/assert"
)

type A struct {
}

func (a A) DoIt() string {
	return "DoIt from A"
}

func NewA() A {
	return A{}
}

type B struct {
	ARef A
}

func NewB(a A) B {
	return B{a}
}

func (b B) DoIt() string {
	return "DoIt from B + " + b.ARef.DoIt()
}

func (b B) EndpointDefinitions() whttp.EndpointDefinitions {
	return whttp.EndpointDefinitions{
		whttp.EndpointDefinition{
			Name:    "B Endpooint",
			Pattern: "/b",
			Method:  "GET",
			Handler: b.EndpointB,
		},
	}
}

func (b B) EndpointB() whttp.Response {
	return whttp.Ok("")
}

func TestA(t *testing.T) {
	assert := assert.New(t)

	type TestAppConfig struct {
		ConfigA     string
		DatabaseURL string `envconfig:"DATABASE_URL"`
	}

	type TestApp struct {
		Config TestAppConfig
	}

	app := application.New(TestApp{})
	app.Provide(NewA)
	app.Provide(NewB)

	app.Invoke(func(config TestAppConfig) int {
		assert.Equal("postgres://weezr:weezr@localhost:5432/weezr?sslmode=disable", config.DatabaseURL)
		return 0
	})

	app.Invoke(func(a A, b B) int {
		assert.Equal("DoIt from A", a.DoIt())
		assert.Equal("DoIt from B + DoIt from A", b.DoIt())
		return 0
	})

	app.ServeHTTP()

}
