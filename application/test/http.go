package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Post(url string, payload interface{}) *http.Request {
	logrus.Infof("POST %s", url)
	requestBody, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err.Error())
	}

	req.Header["Authorization-Test"] = []string{"{\"Id\": \"fake-user\"}"}

	return req
}

func Exec(router *mux.Router, request *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, request)
	return resp
}
