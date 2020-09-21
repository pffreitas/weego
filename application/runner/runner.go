package runner

import (
	"log"

	"github.com/pffreitas/weego/application"
	whttp "github.com/pffreitas/weego/application/http"
)

func Run(app *application.WeegoApplication) {

	app.Invoke(func() int {
		log.Printf("Running: %s", app.Name)
		return 0
	})
}

func ServeHTTP(app *application.WeegoApplication) {
	app.Provide(whttp.NewRouter)

	//app.Invoke(func(router *mux.Router) int {
	//	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", "8000"), router))
	//	return 0
	//})
}
