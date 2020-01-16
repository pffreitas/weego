package application_test

import (
	"testing"

	"github.com/pffreitas/weego/application"
	whttp "github.com/pffreitas/weego/application/http"
	"github.com/pffreitas/weego/application/runner"
	"github.com/stretchr/testify/assert"
)

type A struct {
}

func (a A) DoIt() string {
	return "DoIt from A"
}

type B struct {
	ARef A
}

func (b B) DoIt() string {
	return "DoIt from B + " + b.ARef.DoIt()
}

type C struct {
	A
}

func (c C) DoIt() string {
	return "DoIt from C + " + c.A.DoIt()
}

func (b B) EndpointDefinitions() whttp.EndpointDefinitions {
	return whttp.EndpointDefinitions{
		{
			Name:    "B Endpoint",
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
	assertions := assert.New(t)

	type TestAppConfig struct {
		ConfigA string `envconfig:"CONFIG_A"`
	}

	type DatabaseConfig struct {
		DatabaseURL string `envconfig:"DATABASE_URL"`
	}

	type TestApp struct {
		Config   TestAppConfig
		DbConfig DatabaseConfig
	}

	app := application.New(TestApp{})
	app.Provide(A{})
	app.Provide(B{})
	app.Provide(C{})

	app.Invoke(func(config TestAppConfig, dbConfig DatabaseConfig) int {
		assertions.Equal("AAA", config.ConfigA)
		assertions.Equal("postgres://weezr:weezr@localhost:5432/weezr?sslmode=disable", dbConfig.DatabaseURL)
		return 0
	})

	app.Invoke(func(a A, b B) int {
		assertions.Equal("DoIt from A", a.DoIt())
		assertions.Equal("DoIt from B + DoIt from A", b.DoIt())
		return 0
	})

	app.Invoke(func(c C) int {
		assertions.Equal("DoIt from C + DoIt from A", c.DoIt())
		return 0
	})

	runner.Run(&app)
}
