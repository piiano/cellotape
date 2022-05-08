package router

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router/schema_validator"
	"github.com/piiano/restcontroller/router/utils"
	"reflect"
)

func validateOpenAPIRouter(oa *openapi, flatOperations []operation) error {
	errs := utils.NewErrorsCollector()
	errs.AddIfNotNil(oa.options.HandleAllContentTypes, validateContentTypes(*oa))
	declaredOperation := utils.NewSet[string]()
	for _, flatOp := range flatOperations {
		if !declaredOperation.Add(flatOp.id) {
			errs.AddErrorsIfNotNil(fmt.Errorf("multiple handlers found for operation id %q", flatOp.id))
			continue
		}
		_, found := oa.spec.findSpecOperationByID(flatOp.id)
		if !found {
			errs.AddErrorsIfNotNil(fmt.Errorf("handler recieved for non exising operation id %q is spec - %s", flatOp.id, flatOp.sourcePosition))
		}
		errs.AddErrorsIfNotNil(validateOperation(*oa, flatOp))
		if errs.ErrorOrNil() != nil {
			continue
		}
	}
	errs.AddIfNotNil(oa.options.HandleAllOperations, validateHandleAllOperations(oa, declaredOperation))
	return errs.ErrorOrNil()
}

func validateHandleAllOperations(oa *openapi, declaredOperation utils.Set[string]) error {
	errs := utils.NewErrorsCollector()
	for _, pathItem := range oa.spec.Paths {
		for _, specOp := range pathItem.Operations() {
			if !declaredOperation.Has(specOp.OperationID) {
				errs.AddErrorsIfNotNil(fmt.Errorf("missing handler for operation id %q", specOp.OperationID))
			}
		}
	}
	return errs.ErrorOrNil()
}

func validateContentTypes(oa openapi) error {
	errs := utils.NewErrorsCollector()
	specContentTypes := oa.spec.findSpecContentTypes()
	for _, specContentType := range specContentTypes {
		_, exist := oa.contentTypes[specContentType]
		if !exist {
			errs.AddIfNotNil(ReturnError, fmt.Errorf("content type %s is declared in the spec but has no implementation in the router", specContentType))
		}
	}
	return errs.ErrorOrNil()
}

// check
func validateOperation(oa openapi, operation operation) error {
	errs := utils.NewErrorsCollector()
	specOp, found := oa.spec.findSpecOperationByID(operation.id)
	options := oa.options.OperationValidationOptions(operation.id)

	if !found {
		return fmt.Errorf("operation id %q not found in spec", operation.id)
	}

	for _, chainHandler := range operation.handlers {
		errs.AddIfNotNil(options.ValidateRequestBody, validateRequestBodyType(chainHandler, specOp.RequestBody, oa))
		errs.AddIfNotNil(options.ValidatePathParams, validatePathParamsType(chainHandler, specOp.Parameters, oa))
		errs.AddIfNotNil(options.ValidateQueryParams, validateQueryParamsType(chainHandler, specOp.Parameters, oa))
		errs.AddIfNotNil(options.ValidateResponses, validateResponseTypes(chainHandler, specOp.Operation, oa))
	}
	errs.AddIfNotNil(options.HandleAllOperationResponses, validateHandleAllResponses(operation, specOp))
	return errs.ErrorOrNil()
}

// validateHandleAllResponses checks that every response defined in the spec is handled at least once in the handlers chain
func validateHandleAllResponses(operation operation, specOp specOperation) error {
	errs := utils.NewErrorsCollector()
	handlers := append(operation.handlers, operation.handler)
	responseCodes := utils.NewSet[int](utils.ConcatSlices[int](utils.Map(handlers, func(h handler) []int {
		return utils.Keys(h.responses)
	})...)...)
	for statusStr := range specOp.Responses {
		status, err := parseStatus(statusStr)
		if err != nil {
			errs.AddErrorsIfNotNil(fmt.Errorf("spec declaes an invalid status %s on operation %s", statusStr, operation.id))
			continue
		}
		if responseCodes.Has(status) {
			continue
		}
		errs.AddErrorsIfNotNil(fmt.Errorf("response %d is declared on operation %s but is not declared in the spec - %s", status, operation.id, operation.sourcePosition))
	}
	return errs.ErrorOrNil()
}

// validateRequestBodyType check that a request body type declared on a handler is declared on the spec with a compatible schema.
// a handler does not have to declare and handle the request body defined in the spec, but it can not declare request body which is not defined or incompatible.
func validateRequestBodyType(handler handler, specBody *openapi3.RequestBodyRef, oa openapi) error {
	bodyType := handler.request.requestBody
	validator := schema_validator.NewTypeSchemaValidator(bodyType, openapi3.Schema{}, schema_validator.Options{})
	if bodyType == nilType {
		return nil
	}
	if specBody == nil {
		return fmt.Errorf("handler body type is %s while in spec there is no request body - %s", bodyType, handler.sourcePosition)
	}
	for _, mediaType := range specBody.Value.Content {
		// TODO: allow different media types to fine tune behaviour of validation (for example use non json tags during struct validation for non json mime type)
		if err := validator.WithSchema(*mediaType.Schema.Value).Validate(); err != nil {
			return err
		}
	}
	return nil
}

// validatePathParamsType check that all path params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all path parameters defined in the spec, but it can not declare parameters which are not defined.
func validatePathParamsType(handler handler, specParameters openapi3.Parameters, oa openapi) error {
	return validateParamsType("path", "uri", handler.request.pathParams, specParameters, oa, handler.sourcePosition)
}

// validatePathParamsType check that all query params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all query parameters defined in the spec, but it can not declare parameters which are not defined.
func validateQueryParamsType(handler handler, specParameters openapi3.Parameters, oa openapi) error {
	return validateParamsType("query", "form", handler.request.queryParams, specParameters, oa, handler.sourcePosition)
}

// validateParamsType check that all params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all parameters defined in the spec, but it can not declare parameters which are not defined.
func validateParamsType(in string, tag string, paramsType reflect.Type, specParameters openapi3.Parameters, oa openapi, position sourcePosition) error {
	errs := utils.NewErrorsCollector()
	if paramsType == nilType {
		return nil
	}
	validator := schema_validator.NewTypeSchemaValidator(nilType, openapi3.Schema{}, oa.options.SchemaValidation)
	for _, field := range reflect.VisibleFields(paramsType) {
		name := field.Tag.Get(tag)
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		specParameter := specParameters.GetByInAndName(in, name)
		if specParameter == nil {
			errs.AddErrorsIfNotNil(fmt.Errorf("%s param %q defined in %s pathParams type %s is missing in the spec - %s", in, name, in, paramsType, position))
			continue
		}
		// TODO: schema validator check object schemas with json keys
		errs.AddErrorsIfNotNil(validator.WithType(field.Type).WithSchema(*specParameter.Schema.Value).Validate())
	}
	return errs.ErrorOrNil()
}

// validateResponseTypes check that all responses declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all possible responses defined in the spec, but it can not declare responses which are not defined.
func validateResponseTypes(handler handler, specOperation *openapi3.Operation, oa openapi) error {
	errs := utils.NewErrorsCollector()
	schemaValidator := schema_validator.NewTypeSchemaValidator(reflect.TypeOf(nil), openapi3.Schema{}, oa.options.SchemaValidation)
	for status, response := range handler.responses {
		specResponse := specOperation.Responses.Get(status)
		if specResponse == nil {
			return fmt.Errorf("response %d is declared on an handler for operation %s but is not part of the spec - %s", status, specOperation.OperationID, handler.sourcePosition)
		}
		responseValidator := schemaValidator.WithType(response.responseType)
		for _, mediaType := range specResponse.Value.Content {
			// TODO: allow different media types to fine tune behaviour of validation (for example use non json tags during struct validation for non json mime type)
			errs.AddErrorsIfNotNil(responseValidator.WithSchema(*mediaType.Schema.Value).Validate())
		}
	}
	return errs.ErrorOrNil()
}
