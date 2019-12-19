package runner

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/pffreitas/weego/application"
	whttp "github.com/pffreitas/weego/application/http"
)

func Run(app application.WeegoApplication) {
	fmt.Printf(">>> %+v\n", app)
}

func ServeHTTP(app *application.WeegoApplication) {
	app.Provide(whttp.NewRouter)

	app.Invoke(func(router *mux.Router) int {
		//TODO config port
		// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ""), router))
		return 0
	})
}
