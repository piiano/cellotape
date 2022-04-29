package router

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/schema_validator"
	"github.com/piiano/restcontroller/utils"
	"net/http"
	"reflect"
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
	errs := utils.NewErrorsCollector()
	types := o.getHandlerFunc().types()
	specOp, found := oa.getSpec().findSpecOperationByID(o.id)
	if !found {
		return fmt.Errorf("operation id %q not found in spec", o.id)
	}
	errs.AddIfNotNil(validatePathParamsType(types.pathParams, specPathParameters(specOp), oa))
	errs.AddIfNotNil(validateQueryParamsType(types.queryParams, specQueryParameters(specOp), oa))
	errs.AddIfNotNil(validateRequestBodyType(types.requestBody, specOp.RequestBody, oa.getContentTypes()))
	errs.AddIfNotNil(validateResponseTypes(types.responsesType, specOp.Responses, oa))
	return errs.ErrorOrNil()
}

func specPathParameters(specOp specOperation) openapi3.Parameters {
	return utils.Filter(specOp.Parameters, func(param *openapi3.ParameterRef) bool {
		return param != nil && param.Value != nil && param.Value.In == "path"
	})
}
func specQueryParameters(specOp specOperation) openapi3.Parameters {
	return utils.Filter(specOp.Parameters, func(param *openapi3.ParameterRef) bool {
		return param != nil && param.Value != nil && param.Value.In == "query"
	})
}

func validatePathParamsType(pathParamsType reflect.Type, specPathParameters openapi3.Parameters, oa OpenAPI) error {
	errs := utils.NewErrorsCollector()
	specParameterNames := utils.Map(specPathParameters, func(parameter *openapi3.ParameterRef) string {
		return parameter.Value.Name
	})
	if pathParamsType == nilType && len(specPathParameters) > 0 {
		return fmt.Errorf("path params type %s is incompatible with the spec defines path parameters %q", pathParamsType, specParameterNames)
	}
	if pathParamsType == nilType {
		return nil
	}
	validator := schema_validator.NewTypeSchemaValidator(nilType, openapi3.Schema{}, oa.getOptions().InitializationSchemaValidation)
	declaredParams := make(map[string]bool)
	for _, field := range reflect.VisibleFields(pathParamsType) {
		name := field.Tag.Get("uri")
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		declaredParams[name] = true
		specParameter := specPathParameters.GetByInAndName("path", name)
		if specParameter == nil {
			errs.AddIfNotNil(fmt.Errorf("path param %q defined in path params type %s is missing in the spec", name, pathParamsType))
			continue
		}
		// TODO: schema validator check object schemas with json keys
		errs.AddIfNotNil(validator.WithType(field.Type).WithSchema(*specParameter.Schema.Value).Validate())
	}
	for _, name := range specParameterNames {
		if !declaredParams[name] {
			errs.AddIfNotNil(fmt.Errorf("path param %q defined spec but is missing in path params type %s", name, pathParamsType))
		}
	}
	return errs.ErrorOrNil()
}
func validateQueryParamsType(queryParamsType reflect.Type, specQueryParameters openapi3.Parameters, oa OpenAPI) error {
	errs := utils.NewErrorsCollector()
	specParameterNames := utils.Map(specQueryParameters, func(parameter *openapi3.ParameterRef) string {
		return parameter.Value.Name
	})
	if queryParamsType == nilType && len(specQueryParameters) > 0 {
		return fmt.Errorf("query params type %s is incompatible with the spec defines query parameters %q", queryParamsType, specParameterNames)
	}
	if queryParamsType == nilType {
		return nil
	}
	validator := schema_validator.NewTypeSchemaValidator(nilType, openapi3.Schema{}, oa.getOptions().InitializationSchemaValidation)
	declaredParams := make(map[string]bool)
	for _, field := range reflect.VisibleFields(queryParamsType) {
		name := field.Tag.Get("form")
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		declaredParams[name] = true
		specParameter := specQueryParameters.GetByInAndName("query", name)
		if specParameter == nil {
			errs.AddIfNotNil(fmt.Errorf("query param %q defined in query params type %s is missing in the spec", name, queryParamsType))
			continue
		}
		// TODO: schema validator check object schemas with json keys
		errs.AddIfNotNil(validator.WithType(field.Type).WithSchema(*specParameter.Schema.Value).Validate())
	}
	for _, name := range specParameterNames {
		if !declaredParams[name] {
			errs.AddIfNotNil(fmt.Errorf("query param %q defined spec but is missing in query params type %s", name, queryParamsType))
		}
	}
	return errs.ErrorOrNil()
}

func validateResponseTypes(responsesType reflect.Type, operationResponses openapi3.Responses, oa OpenAPI) error {
	errs := utils.NewErrorsCollector()
	responsesMap, err := extractResponses(responsesType)
	if err != nil {
		return err
	}
	supportedContentTypes := oa.getContentTypes()
	schemaValidator := schema_validator.NewTypeSchemaValidator(
		reflect.TypeOf(nil),
		openapi3.Schema{},
		oa.getOptions().InitializationSchemaValidation,
	)
	declaredResponses := make(map[int]bool)
	for statusStr, specResponse := range operationResponses {
		status, err := parseStatus(statusStr)
		if errs.AddIfNotNil(err) {
			continue
		}
		responseType, found := responsesMap[status]
		if !found {
			errs.AddIfNotNil(fmt.Errorf("spec response for status %d is not declared in the responses type %s", status, responsesType))
			continue
		}
		declaredResponses[status] = true
		responseValidator := schemaValidator.WithType(responseType.responseType)
		for mime, mediaType := range specResponse.Value.Content {
			_, found := supportedContentTypes[mime]
			if !found {
				errs.AddIfNotNil(fmt.Errorf("response content type with mime value %q is missing", mime))
				continue
			}
			errs.AddIfNotNil(responseValidator.WithSchema(*mediaType.Schema.Value).Validate())
		}
	}
	for status, _ := range responsesMap {
		if declaredResponses[status] {
			continue
		}
		errs.AddIfNotNil(fmt.Errorf("response status %d of responses type %s is not declared in the spec", status, responsesType))
	}
	return errs.ErrorOrNil()
}
