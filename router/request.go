package router

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/schema_validator"
	"net/http"
	"reflect"
)

type Request[B, P, Q any] struct {
	Context     context.Context
	Body        B
	PathParams  P
	QueryParams Q
	Headers     http.Header
}

type RequestTypes struct {
	requestBody reflect.Type
	pathParams  reflect.Type
	queryParams reflect.Type
}

type Nil *uintptr

var nilValue Nil
var nilType = reflect.TypeOf(nilValue)

func validateRequestBodyType(bodyType reflect.Type, specBody *openapi3.RequestBodyRef, contentTypes ContentTypes) error {
	validator := schema_validator.NewTypeSchemaValidator(bodyType, openapi3.Schema{}, schema_validator.Options{})
	if specBody == nil && bodyType == nilType {
		return nil
	}
	if specBody == nil {
		return fmt.Errorf("operation handler body type is %s while in spec there is no request body", bodyType)
	}
	if bodyType == nilType {
		return fmt.Errorf("operation handler body type is %s while the spec has a request body", bodyType)
	}
	for mime, mediaType := range specBody.Value.Content {
		_, found := contentTypes[mime]
		if !found {
			return fmt.Errorf("missing handler for content type with mime value %q defined in spec", mime)
		}
		if err := validator.WithSchema(*mediaType.Schema.Value).Validate(); err != nil {
			return err
		}
	}
	return nil
}
