package http

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

func bindBody(payload reflect.Value, r *http.Request) {
	bodyField := payload.FieldByName("Body")

	if bodyField.IsValid() {
		body := reflect.New(bodyField.Type()).Interface()

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err := decoder.Decode(body)
		if err != nil {
			//TODO handler err
		}

		bodyField.Set(reflect.ValueOf(body).Elem())
	}
}

func parseRequestParams(r *http.Request) map[string]string {
	params := map[string]string{}
	for k, v := range mux.Vars(r) {
		params[k] = v
	}

	for k, v := range r.URL.Query() {
		params[k] = v[0]
	}

	return params
}

func bindParams(payload reflect.Value, params map[string]string) {
	paramsField := payload.FieldByName("Params")

	if paramsField.IsValid() {
		paramsStruct := reflect.New(paramsField.Type()).Elem()

		for k, v := range params {
			f := paramsStruct.FieldByName(strings.Title(k))
			if f.IsValid() {
				f.Set(reflect.ValueOf(v))
			}
		}

		paramsField.Set(paramsStruct)
	}
}
