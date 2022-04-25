package router

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPISpec openapi3.T

func NewSpecFromFile(path string) (OpenAPISpec, error) {
	spec, err := openapi3.NewLoader().LoadFromFile(path)
	if err != nil {
		return OpenAPISpec{}, err
	}
	return OpenAPISpec(*spec), nil
}

func NewSpecFromData(data []byte) (OpenAPISpec, error) {
	spec, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		return OpenAPISpec{}, err
	}
	return OpenAPISpec(*spec), nil
}

func NewSpec() OpenAPISpec {
	spec, _ := NewSpecFromData([]byte("{}"))
	return spec
}
