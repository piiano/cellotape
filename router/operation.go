package router

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/schema_validator"
	"math/bits"
	"net/http"
	"reflect"
	"strconv"
)

type Operation interface {
	getId() string
	getHandlerFunc() OperationHandler
	getHandlers() []http.HandlerFunc
	validateHandlerTypes(oa OpenAPI) error
}
type operation struct {
	id               string
	handlers         []http.HandlerFunc
	operationHandler OperationHandler
}

func (o operation) getId() string                    { return o.id }
func (o operation) getHandlerFunc() OperationHandler { return o.operationHandler }
func (o operation) getHandlers() []http.HandlerFunc  { return o.handlers }

func (o operation) validateHandlerTypes(oa OpenAPI) error {
	types, err := o.getHandlerFunc().types()
	if err != nil {
		return err
	}
	schemaValidator := schema_validator.NewTypeSchemaValidator(reflect.TypeOf(nil), openapi3.Schema{}, oa.getOptions().InitializationSchemaValidation)
	specOp, found := oa.getSpec().findSpecOperationByID(o.id)
	if !found {
		return fmt.Errorf("operation id %q not found in spec", o.id)
	}
	if err := ValidateRequestBody(types.RequestBody, specOp.RequestBody, oa.getContentTypes()); err != nil {
		return err
	}

	for _, param := range specOp.Parameters {
		// TODO: validate path params schema
		//paramSchema := param.Value.Schema.Value
		if param.Value.In == "path" {
			//fmt.Println(paramSchema, types.PathParams)
		}
		// TODO: validate query params schema
		if param.Value.In == "query" {
			//fmt.Println(paramSchema, types.QueryParams)
		}
	}

	responseValidator := schemaValidator.WithType(types.ResponseType)
	for statusStr, response := range specOp.Responses {
		for mime, mediaType := range response.Value.Content {
			_, found = oa.getContentTypes()[mime]
			if !found {
				return fmt.Errorf("content type with mime value %q not found in spec", mime)
			}
			if err = responseValidator.WithSchema(*mediaType.Schema.Value).Validate(); err != nil {
				return err
			}
			status, err := strconv.ParseInt(statusStr, 10, bits.UintSize)
			if err != nil {
				return err
			}
			// spec responses can be for both error and successful cases.
			// there might be more than one successful and one error response.
			// we only have one success type and one error type.
			// we need to think what is the validation strategy here.
			// TODO: validate responses schema
			if status < 300 && status >= 200 {
				continue
			}
			// TODO: validate http error schema

		}
	}
	return nil
}

//
//// how to represent multiple response options in a way their type can be extracted?
//type responseFunc[B, P, Q, R any, H interface {
//	~func(request Request[B, P, Q]) (R, error)
//}] func(int, contentType, value R) HttpResponse
//
////type responseFunc[R any] func(int, contentType, value R) HttpResponse
//type operationHandlerFuncV2[B, P, Q, R any] func(request Request[B, P, Q], send responseFunc[R]) error

/*
# Response Validation Flow

type Response[T any] func(status, value T) HttpResponse


WithOperation(
	NewOperation(func... (...) HttpResponse).
	WithResponse(200, any?) // Content Type?
	WithResponse(410, any?)
)
*/
//
//type HttpResponseV2 interface {
//	getStatus() int
//	getContentType() string
//	getBodyType() (reflect.Type, error)
//	bytes() ([]byte, error)
//}
//
//type httpResponseV2[T any] struct {
//	status      int
//	contentType string
//	body        T
//}
//
//func NewHTTPResponseV2[T any](status int, contentType string) HttpResponseV2 {
//	return httpResponse[T]{
//		status:      status,
//		contentType: contentType,
//	}
//}
