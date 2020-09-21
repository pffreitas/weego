package application

import (
	"reflect"

	"github.com/pffreitas/weego/application/http"
)

type WeegoApplication struct {
	Name      string
	container *container
}

func New(app interface{}) WeegoApplication {
	container := newContainer()

	configObjects, err := processConfig(app)
	if err == nil {
		for _, configObject := range configObjects {
			container.injectConfig(configObject)
		}
	}

	weegoApplication := container.invoke(newWeegoApplication).(WeegoApplication)

	appVal := reflect.ValueOf(app)
	weegoApplication.Name = appVal.Type().Name()

	return weegoApplication
}

func newWeegoApplication(container *container) WeegoApplication {
	return WeegoApplication{"", container}
}

func (wa *WeegoApplication) Provide(constructor interface{}) {
	wa.container.provide(constructor)
}

func (wa *WeegoApplication) Invoke(fn interface{}) {
	wa.container.invoke(fn)
}

func (wa *WeegoApplication) Use(middlewareFns ...http.Middleware) {
	for _, fn := range middlewareFns {
		wa.container.provide(fn)
	}
}
