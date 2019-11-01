package http

import (
	"net/http"
	"reflect"
)

func getHandlerFuncParamTypes(handlerFunc interface{}) []reflect.Type {
	handlerFuncType := reflect.TypeOf(handlerFunc)
	numArgs := handlerFuncType.NumIn()

	handlerFuncParamTypes := make([]reflect.Type, 0, numArgs)
	for i := 0; i < numArgs; i++ {
		p := handlerFuncType.In(i)
		handlerFuncParamTypes = append(handlerFuncParamTypes, p)
	}

	return handlerFuncParamTypes
}

func createHandler(handlerFunc interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		handlerFuncParamTypes := getHandlerFuncParamTypes(handlerFunc)
		handlerFuncArgs := make([]reflect.Value, 0, len(handlerFuncParamTypes))

		if len(handlerFuncParamTypes) > 0 {
			payload := reflect.New(handlerFuncParamTypes[0]).Elem()
			bindBody(payload, request)
			bindParams(payload, parseRequestParams(request))
			handlerFuncArgs = append(handlerFuncArgs, payload)
		}

		response := reflect.ValueOf(handlerFunc).Call(handlerFuncArgs)
		whttpResponse := response[0].Interface().(Response)

		writeJSON(w, whttpResponse.Body, whttpResponse.Code)
	}
}