package http

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

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
			logrus.WithError(err).Error("failed to decode body payload")
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

var (
	converters = map[string]interface{}{
		"int64": int64Converter,
	}
)

func int64Converter(param string) (int64, error) {
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		logrus.WithError(err).Errorf("failed to parse to int64: ", param)
		return 0, err
	}

	return val, err
}

func bindParams(payload reflect.Value, params map[string]string) {
	paramsField := payload.FieldByName("Params")

	if paramsField.IsValid() {
		paramsStruct := reflect.New(paramsField.Type()).Elem()

		for k, v := range params {
			f := paramsStruct.FieldByName(strings.Title(k))

			if f.IsValid() {
				converter := converters[f.Type().Name()]
				if converter != nil {
					out := reflect.ValueOf(converter).Call([]reflect.Value{reflect.ValueOf(v)})
					if out[1].IsNil() {
						f.Set(out[0])
					}
				} else {
					f.Set(reflect.ValueOf(v))
				}
			}
		}

		paramsField.Set(paramsStruct)
	}
}
