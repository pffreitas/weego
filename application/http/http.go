package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/robbert229/jwt"

	"github.com/gorilla/mux"
)

type Middleware interface {
	Handle(w http.ResponseWriter, request *http.Request) bool
}

type CorsMiddleware struct {
}

func (m CorsMiddleware) Handle(w http.ResponseWriter, request *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	if "OPTIONS" == request.Method {
		writeJSON(w, "", 200)
		return false
	}

	return true
}

type AuthMiddleware struct {
}

func (m AuthMiddleware) Handle(w http.ResponseWriter, request *http.Request) bool {
	authorizationHeader := request.Header["Authorization"]

	logrus.WithField("authorization", authorizationHeader).Info("authorizing")

	if len(authorizationHeader) == 1 {
		authToken := authorizationHeader[0]

		if strings.HasPrefix(authToken, "Bearer ") {
			jwtToken := strings.Replace(authToken, "Bearer ", "", 1)

			hs256 := jwt.HmacSha256("config.JWTSecret")

			if hs256.Validate(jwtToken) != nil {
				logrus.Errorf("Invalid JWT token %s", jwtToken)
				writeJSON(w, "Unauthorized", 401)
				return false
			}

			if claims, err := hs256.Decode(jwtToken); err == nil {
				if sub, err := claims.Get("sub"); err == nil {
					w.Header().Set("Username", fmt.Sprintf("%v", sub))
					return true
				}
			}
		}
	}

	writeJSON(w, "Unauthorized", 401)
	logrus.Errorf("Failed to process authorization: %v", authorizationHeader)
	return false
}

func NewRouter(endpoints []EndpointProvider, middlewareFns []Middleware) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, e := range endpoints {
		fmt.Printf("- Registering %+v endpoints \n", reflect.TypeOf(e))

		for _, ed := range e.EndpointDefinitions() {
			fmt.Printf("\t %v \n", ed)

			router.
				Methods(ed.Method).
				Path(ed.Pattern).
				Name(ed.Name).
				Handler(createHandler(middlewareFns, ed.Handler))
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
