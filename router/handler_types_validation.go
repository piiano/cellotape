package router

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/schema_validator"
	"github.com/piiano/restcontroller/utils"
	"reflect"
)

func validateRequestBodyType(bodyType reflect.Type, specBody *openapi3.RequestBodyRef, oa openapi) error {
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
		_, found := oa.contentTypes[mime]
		if !found {
			return fmt.Errorf("missing handler for content type with mime value %q defined in spec", mime)
		}
		if err := validator.WithSchema(*mediaType.Schema.Value).Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (o operation) validateHandlerTypes(oa openapi) error {
	errs := utils.NewErrorsCollector()
	types := o.operationHandler.types()
	specOp, found := oa.spec.findSpecOperationByID(o.id)
	if !found {
		return fmt.Errorf("operation id %q not found in spec", o.id)
	}
	errs.AddIfNotNil(validatePathParamsType(types.pathParams, specOp.Parameters, oa))
	errs.AddIfNotNil(validateQueryParamsType(types.queryParams, specOp.Parameters, oa))
	errs.AddIfNotNil(validateRequestBodyType(types.requestBody, specOp.RequestBody, oa))
	errs.AddIfNotNil(validateResponseTypes(types.responsesType, specOp.Responses, oa))
	return errs.ErrorOrNil()
}

func validatePathParamsType(pathParamsType reflect.Type, specParameters openapi3.Parameters, oa openapi) error {
	return validateParamsType("path", "uri", pathParamsType, specParameters, oa)
}

func validateQueryParamsType(queryParamsType reflect.Type, specParameters openapi3.Parameters, oa openapi) error {
	return validateParamsType("query", "form", queryParamsType, specParameters, oa)
}

func validateParamsType(in string, tag string, paramsType reflect.Type, specParameters openapi3.Parameters, oa openapi) error {
	specParameters = utils.Filter(specParameters, func(param *openapi3.ParameterRef) bool {
		return param != nil && param.Value != nil && param.Value.In == in
	})
	errs := utils.NewErrorsCollector()
	specParameterNames := utils.Map(specParameters, func(parameter *openapi3.ParameterRef) string {
		return parameter.Value.Name
	})
	if paramsType == nilType && len(specParameters) > 0 {
		return fmt.Errorf("%s params type %s is incompatible with the spec defines %s parameters %v", in, paramsType, in, specParameterNames)
	}
	if paramsType == nilType {
		return nil
	}
	validator := schema_validator.NewTypeSchemaValidator(nilType, openapi3.Schema{}, oa.options.InitializationSchemaValidation)
	declaredParams := make(map[string]bool)
	for _, field := range reflect.VisibleFields(paramsType) {
		name := field.Tag.Get(tag)
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		declaredParams[name] = true
		specParameter := specParameters.GetByInAndName(in, name)
		if specParameter == nil {
			errs.AddIfNotNil(fmt.Errorf("%s param %q defined in %s params type %s is missing in the spec", in, name, in, paramsType))
			continue
		}
		// TODO: schema validator check object schemas with json keys
		errs.AddIfNotNil(validator.WithType(field.Type).WithSchema(*specParameter.Schema.Value).Validate())
	}
	for _, name := range specParameterNames {
		if !declaredParams[name] {
			errs.AddIfNotNil(fmt.Errorf("%s param %q defined spec but is missing in %s params type %s", in, name, in, paramsType))
		}
	}
	return errs.ErrorOrNil()
}

func validateResponseTypes(responsesType reflect.Type, operationResponses openapi3.Responses, oa openapi) error {
	errs := utils.NewErrorsCollector()
	responsesMap, err := extractResponses(responsesType)
	if err != nil {
		return err
	}
	supportedContentTypes := oa.contentTypes
	schemaValidator := schema_validator.NewTypeSchemaValidator(
		reflect.TypeOf(nil),
		openapi3.Schema{},
		oa.options.InitializationSchemaValidation,
	)
	declaredResponses := make(map[int]bool)
	for statusStr, specResponse := range operationResponses {
		status, err := parseStatus(statusStr)
		if errs.AddIfNotNil(err) {
			continue
		}
		responseType, found := responsesMap[status]
		if !found {
			errs.AddIfNotNil(fmt.Errorf("spec httpResponse for status %d is not declared in the responses type %s", status, responsesType))
			continue
		}
		declaredResponses[status] = true
		responseValidator := schemaValidator.WithType(responseType.responseType)
		for mime, mediaType := range specResponse.Value.Content {
			_, found := supportedContentTypes[mime]
			if !found {
				errs.AddIfNotNil(fmt.Errorf("httpResponse content type with mime value %q is missing", mime))
				continue
			}
			errs.AddIfNotNil(responseValidator.WithSchema(*mediaType.Schema.Value).Validate())
		}
	}
	for status := range responsesMap {
		if declaredResponses[status] {
			continue
		}
		errs.AddIfNotNil(fmt.Errorf("httpResponse status %d of responses type %s is not declared in the spec", status, responsesType))
	}
	return errs.ErrorOrNil()
}
