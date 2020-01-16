package http

import (
	"fmt"
	"net/http"
)

type Response struct {
	Code int
	Body interface{}
}

func Ok(body interface{}) Response {
	return Response{
		http.StatusOK,
		body,
	}
}

func Created(body interface{}) Response {
	return Response{
		http.StatusCreated,
		body,
	}
}

func ServerError(err error) Response {
	return Response{
		http.StatusInternalServerError,
		err,
	}
}

type EndpointDefinition struct {
	Name    string
	Pattern string
	Method  string
	Handler interface{}
}

func (ed EndpointDefinition) String() string {
	return fmt.Sprintf("%-20s: %-6s %s", ed.Name, fmt.Sprintf("[%s]", ed.Method), ed.Pattern)
}

type EndpointDefinitions []EndpointDefinition

type EndpointProvider interface {
	EndpointDefinitions() EndpointDefinitions
}
