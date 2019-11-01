package application_test

import (
	"testing"

	"github.com/pffreitas/weego/application"
	"github.com/stretchr/testify/assert"
)

type A struct {
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

	app.Run()

	app.Invoke(func(config TestAppConfig) int {
		assert.Equal("postgres://weezr:weezr@localhost:5432/weezr?sslmode=disable", config.DatabaseURL)
		return 0
	})
}
