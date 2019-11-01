package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// NewRouter ..
func NewRouter(endpoints []EndpointProvider) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, e := range endpoints {
		fmt.Printf("- Registering %+v endpoints \n", reflect.TypeOf(e))

		for _, ed := range e.EndpointDefinitions() {
			fmt.Printf("\t %v \n", ed)

			router.
				Methods(ed.Method).
				Path(ed.Pattern).
				Name(ed.Name).
				Handler(createHandler(ed.Handler))
		}
	}

	return router
}

func writeJSON(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			// TODO handler err
		}

		w.Write(bytes)
	}
}
