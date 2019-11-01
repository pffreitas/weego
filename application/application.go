package application

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	whttp "github.com/pffreitas/weego/application/http"
)

// WeegoApplication .
type WeegoApplication struct {
	container *container
	router    *mux.Router
}

// New .
func New(app interface{}) WeegoApplication {
	container := newContainer()

	config, err := processConfig(app)
	if err == nil {
		container.injectConfig(config)
	}

	container.provide(whttp.NewRouter)

	return container.invoke(newWeegoApplication).(WeegoApplication)
}

func newWeegoApplication(container *container, router *mux.Router) WeegoApplication {
	return WeegoApplication{container, router}
}

// ServeHTTP .
func (wa *WeegoApplication) ServeHTTP() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", "9004"), wa.router))
}

// Provide .
func (wa *WeegoApplication) Provide(constructor interface{}) {
	wa.container.provide(constructor)
}

// Invoke .
func (wa *WeegoApplication) Invoke(fn interface{}) {
	wa.container.invoke(fn)
}
