package application_test

import (
	"fmt"
	"testing"

	"github.com/pffreitas/weego/application/test"

	"github.com/gorilla/mux"

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

func (a A) EndpointDefinitions() whttp.EndpointDefinitions {
	return whttp.EndpointDefinitions{
		{
			Name:    "Save",
			Pattern: "/foo",
			Method:  "POST",
			Handler: a.save,
		},
	}
}

type SavePayload struct {
	Params struct {
		Limit    int64
		Offset   int64
		Foo      string
		Username string
	}
	Body struct {
		ParamA float64
		ParamB string
		ParamC int64
	}
}

func (a A) save(payload SavePayload) whttp.Response {
	return whttp.Ok(payload)
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

func TestB(t *testing.T) {

	type TestApp struct {
	}

	app := application.New(TestApp{})
	app.Use(whttp.AuthMiddleware{}, whttp.CorsMiddleware{})
	app.Provide(A{})

	runner.ServeHTTP(&app)

	s := map[string]interface{}{
		"ParamA": 2.2,
		"ParamB": "B",
		"ParamC": int64(1),
	}

	app.Invoke(func(router *mux.Router) int {
		response := test.Exec(router, test.Post("/foo?Limit=10&offset=20&foo=bar", s))

		fmt.Printf("\n >> %d \n", response.Code)
		fmt.Printf("\n >> %s \n", response.Body.String())

		return 0
	})

}
