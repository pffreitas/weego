package http

import (
	"fmt"
	"net/http"
)

// Response .
type Response struct {
	Code int
	Body interface{}
}

// Ok .
func Ok(body interface{}) Response {
	return Response{
		http.StatusOK,
		body,
	}
}

// Created .
func Created(body interface{}) Response {
	return Response{
		http.StatusCreated,
		body,
	}
}

// ServerError .
func ServerError(err error) Response {
	return Response{
		http.StatusInternalServerError,
		err,
	}
}

// EndpointDefinition .
type EndpointDefinition struct {
	Name    string
	Pattern string
	Method  string
	Handler interface{}
}

func (ed EndpointDefinition) String() string {
	return fmt.Sprintf("%-20s: %-6s %s", ed.Name, fmt.Sprintf("[%s]", ed.Method), ed.Pattern)
}

// EndpointDefinitions .
type EndpointDefinitions []EndpointDefinition

// EndpointProvider .
type EndpointProvider interface {
	EndpointDefinitions() EndpointDefinitions
}
