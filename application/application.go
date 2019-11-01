package application

import (
	"fmt"

	"github.com/gorilla/mux"
	whttp "github.com/pffreitas/weego/application/http"
)

// WeegoApplication .
type WeegoApplication struct {
	container *container
}

// New .
func New(app interface{}) WeegoApplication {
	container := newContainer()

	config, err := processConfig(app)
	if err == nil {
		container.injectConfig(config)
	}

	return container.invoke(newWeegoApplication).(WeegoApplication)
}

func newWeegoApplication(container *container) WeegoApplication {
	return WeegoApplication{container}
}

// ServeHTTP .
func (wa *WeegoApplication) ServeHTTP() {
	wa.container.provide(whttp.NewRouter)

	wa.Invoke(func(router *mux.Router) int {
		fmt.Printf("%+v \n", router)
		//TODO config port
		// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ""), router))
		return 0
	})
}

// Provide .
func (wa *WeegoApplication) Provide(constructor interface{}) {
	wa.container.provide(constructor)
}

// Invoke .
func (wa *WeegoApplication) Invoke(fn interface{}) {
	wa.container.invoke(fn)
}
