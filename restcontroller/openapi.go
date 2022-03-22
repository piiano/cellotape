package restcontroller

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"reflect"
)

type OperationOptions struct {
	Errors openapi3.Responses
}

func (fn ControllerFn[B, P, Q, R]) OpenAPIOperation(ID string, options *OperationOptions) (*openapi3.Operation, error) {
	generator := openapi3gen.NewGenerator()
	controllerTypes := fn.TypeInfo()
	operation := openapi3.NewOperation()
	operation.OperationID = ID
	if err := appendResponses(generator, operation, controllerTypes.SuccessfulResponseBody, &options.Errors); err != nil {
		return nil, err
	}
	if err := appendRequestBody(generator, operation, controllerTypes.RequestBody); err != nil {
		return nil, err
	}
	if err := appendParams(generator, operation, controllerTypes.PathParams, openapi3.NewPathParameter); err != nil {
		return nil, err
	}
	if err := appendParams(generator, operation, controllerTypes.QueryParams, openapi3.NewQueryParameter); err != nil {
		return nil, err
	}
	return operation, nil
}

// Append request body to operation based on body type
func appendRequestBody(generator *openapi3gen.Generator,
	operation *openapi3.Operation,
	requestBodyType reflect.Type) error {
	if requestBodyType.Kind() == reflect.Struct {
		requestBody := openapi3.NewRequestBody()
		schema, err := generator.GenerateSchemaRef(requestBodyType)
		if err != nil {
			return err
		}
		requestBody.Content = openapi3.NewContentWithJSONSchema(schema.Value)
		operation.RequestBody = &openapi3.RequestBodyRef{Value: requestBody}
	}
	return nil
}

// Append the type of 200 response to responses map
func appendResponses(generator *openapi3gen.Generator,
	operation *openapi3.Operation,
	successfulResponseType reflect.Type,
	errors *openapi3.Responses) error {
	operation.Responses = openapi3.NewResponses()
	successfulResponseSchema, err := generator.GenerateSchemaRef(successfulResponseType)
	if err != nil {
		return err
	}
	successfulResponse := openapi3.NewResponse()
	successfulResponse.Content = openapi3.NewContentWithJSONSchema(successfulResponseSchema.Value)
	operation.Responses["200"] = &openapi3.ResponseRef{Value: successfulResponse}
	for errorCode, errorResponse := range *errors {
		operation.Responses[errorCode] = errorResponse
	}
	return nil
}

// Append parameter to operations using struct types to represent the params.
func appendParams(
	generator *openapi3gen.Generator,
	operation *openapi3.Operation,
	paramsType reflect.Type,
	newParam func(name string) *openapi3.Parameter) error {
	if operation.Parameters == nil {
		operation.Parameters = openapi3.NewParameters()
	}
	if paramsType.Kind() != reflect.Struct {
		return nil
	}
	schema, err := generator.GenerateSchemaRef(paramsType)
	if err != nil {
		return err
	}
	for name, property := range schema.Value.Properties {
		parameter := newParam(name)
		parameter.Schema = property
		operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{Value: parameter})
	}
	return nil
}
