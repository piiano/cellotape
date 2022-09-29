package router

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/utils"
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

// Operations returns all operations declared in the spec as a map of operation ID to SpecOperation.
func (s *OpenAPISpec) Operations() map[string]SpecOperation {
	operations := make(map[string]SpecOperation, 0)
	for path, pathItem := range s.Paths {
		for method, specOp := range pathItem.Operations() {
			operations[specOp.OperationID] = SpecOperation{
				Path:      path,
				Method:    method,
				Operation: specOp,
			}
		}
	}
	return operations
}

func (s *OpenAPISpec) findSpecOperationByID(id string) (SpecOperation, bool) {
	for path, pathItem := range s.Paths {
		for method, specOp := range pathItem.Operations() {
			if specOp.OperationID == id {
				return SpecOperation{Path: path, Method: method, Operation: specOp}, true
			}
		}
	}
	return SpecOperation{}, false
}

// findSpecContentTypes find all content types declared in the spec for both request body and responses
func (s *OpenAPISpec) findSpecContentTypes(excludeOperations utils.Set[string]) []string {
	contentTypes := make([]string, 0)
	for _, pathItem := range s.Paths {
		for _, specOp := range pathItem.Operations() {
			if excludeOperations.Has(specOp.OperationID) {
				continue
			}
			if specOp.RequestBody != nil && specOp.RequestBody.Value != nil {
				for contentType := range specOp.RequestBody.Value.Content {
					contentTypes = append(contentTypes, contentType)
				}
			}
			for _, response := range specOp.Responses {
				if response.Value == nil {
					continue
				}
				for contentType := range response.Value.Content {
					contentTypes = append(contentTypes, contentType)
				}
			}
		}
	}
	return contentTypes
}

// SpecOperation represent the operation information described in the spec with path and method information.
type SpecOperation struct {
	Path   string
	Method string
	*openapi3.Operation
}
