package router

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router/schema_validator"
	"github.com/piiano/restcontroller/router/utils"
	"reflect"
)

const (
	pathParamInValue    = "path"
	pathParamFieldTag   = "uri"
	queryParamInValue   = "query"
	queryParamFieldTag  = "form"
	ignoreFieldTagValue = "-"
)

// validateOpenAPIRouter validates the entire OpenAPI Router structure built with the builder with the spec.
// This takes into account various options defined and print to the logs relevant errors and warning based on the defined log level.
func validateOpenAPIRouter(oa *openapi, flatOperations []operation) error {
	l := oa.logger()
	l.ErrorIfNotNil(validateContentTypes(*oa))
	declaredOperation := utils.NewSet[string]()
	for _, flatOp := range flatOperations {
		if !declaredOperation.Add(flatOp.id) {
			// multiple handlers for the same operation is always an error
			l.Errorf(multipleHandlersFoundForOperationId(flatOp.id))
		}
		l.ErrorIfNotNil(validateOperation(*oa, flatOp))
	}
	l.ErrorIfNotNil(validateMustHandleAllOperations(oa, declaredOperation))
	return l.MustHaveNoErrorsf(failedValidatingTheRouterWithTheSpec(l.Warnings(), l.Errors()))
}

// validateMustHandleAllOperations checks that all operations defined in the spec have an implementation on the router.
func validateMustHandleAllOperations(oa *openapi, declaredOperation utils.Set[string]) error {
	l := oa.logger()
	for _, pathItem := range oa.spec.Paths {
		for _, specOp := range pathItem.Operations() {
			if !declaredOperation.Has(specOp.OperationID) {
				l.Logf(utils.LogLevel(oa.options.MustHandleAllOperations), missingHandlerForOperationId(specOp.OperationID))
			}
		}
	}
	return l.MustHaveNoErrorsf(notImplementedSpecOperations(l.Errors()))
}

// validateContentTypes checks that all content types defined in the spec for request or responses have an implementation on the router.
func validateContentTypes(oa openapi) error {
	log := oa.logger()
	level := utils.LogLevel(oa.options.HandleAllContentTypes)
	specContentTypes := oa.spec.findSpecContentTypes()
	for _, specContentType := range specContentTypes {
		_, exist := oa.contentTypes[specContentType]
		if !exist {
			log.Logf(level, missingContentTypeImplementation(specContentType))
		}
	}
	return log.MustHaveNoErrors()
}

// validateOperation perform a validation for an operation and its handlers chain for compliance with the spec.
func validateOperation(oa openapi, operation operation) error {
	l := oa.logger()
	specOp, found := oa.spec.findSpecOperationByID(operation.id)
	options := oa.options.operationValidationOptions(operation.id)
	if !found {
		return fmt.Errorf(handlerForNonExistingSpecOperation(operation.id, operation.sourcePosition))
	}
	for _, chainHandler := range append(operation.handlers, operation.handler) {
		l.ErrorIfNotNil(validateRequestBodyType(oa, options.ValidateRequestBody, chainHandler, specOp.RequestBody, operation.id))
		l.ErrorIfNotNil(validatePathParamsType(oa, options.ValidatePathParams, chainHandler, specOp.Parameters, operation.id))
		l.ErrorIfNotNil(validateQueryParamsType(oa, options.ValidateQueryParams, chainHandler, specOp.Parameters, operation.id))
		l.ErrorIfNotNil(validateResponseTypes(oa, options.ValidateResponses, chainHandler, specOp.Operation, operation.id))
	}
	l.ErrorIfNotNil(validateHandleAllResponses(oa, options.HandleAllOperationResponses, operation, specOp))
	return l.MustHaveNoErrorsf("operation %q has incompatibility with the spec (%d errors, %d warnings)", operation.id, l.Errors(), l.Warnings())
}

// validateHandleAllResponses checks that every response defined in the spec is handled at least once in the handlers chain
func validateHandleAllResponses(oa openapi, behaviour Behaviour, operation operation, specOp specOperation) error {
	l := oa.logger()
	level := utils.LogLevel(behaviour)
	handlers := append(operation.handlers, operation.handler)
	responseCodes := utils.NewSet[int](utils.ConcatSlices[int](utils.Map(handlers, func(h handler) []int {
		return utils.Keys(h.responses)
	})...)...)
	for statusStr := range specOp.Responses {
		status, err := parseStatus(statusStr)
		if err != nil {
			l.Logf(level, invalidStatusInSpecResponses(statusStr, operation.id))
			continue
		}
		if !responseCodes.Has(status) {

			l.Logf(level, handlerDefinesResponseThatIsMissingInSpec(status, operation.id))
		}
	}

	return l.MustHaveNoErrorsf(unimplementedResponsesForOperation(len(specOp.Responses)-len(responseCodes), operation.id))
}

// validateRequestBodyType check that a request body type declared on a handler is declared on the spec with a compatible schema.
// a handler does not have to declare and handle the request body defined in the spec, but it can not declare request body which is not defined or incompatible.
func validateRequestBodyType(oa openapi, behaviour Behaviour, handler handler, specBody *openapi3.RequestBodyRef, operationID string) error {
	l := oa.logger()
	level := utils.LogLevel(behaviour)
	bodyType := handler.request.requestBody
	if bodyType == nilType {
		return nil
	}
	validator := schema_validator.NewTypeSchemaValidator(l, level, bodyType, openapi3.Schema{})
	if specBody == nil {
		l.Logf(level, handlerDefinesRequestBodyWhenNoRequestBodyInSpec(operationID))
		return l.MustHaveNoErrors()
	}
	for _, mediaType := range specBody.Value.Content {
		// TODO: allow different media types to fine tune behaviour of validation (for example use non json tags during struct validation for non json mime type)
		l.LogIfNotNil(level, validator.WithSchema(*mediaType.Schema.Value).Validate())
	}
	return l.MustHaveNoErrors()
}

// validatePathParamsType check that all pathParamInValue params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all pathParamInValue parameters defined in the spec, but it can not declare parameters which are not defined.
func validatePathParamsType(oa openapi, behaviour Behaviour, handler handler, specParameters openapi3.Parameters, operationId string) error {
	return validateParamsType(oa, behaviour, pathParamInValue, pathParamFieldTag, handler.request.pathParams, specParameters, operationId)
}

// validatePathParamsType check that all queryParamInValue params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all queryParamInValue parameters defined in the spec, but it can not declare parameters which are not defined.
func validateQueryParamsType(oa openapi, behaviour Behaviour, handler handler, specParameters openapi3.Parameters, operationId string) error {
	return validateParamsType(oa, behaviour, queryParamInValue, queryParamFieldTag, handler.request.queryParams, specParameters, operationId)
}

// validateParamsType check that all params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all parameters defined in the spec, but it can not declare parameters which are not defined.
func validateParamsType(oa openapi, behaviour Behaviour, in string, tag string, paramsType reflect.Type, specParameters openapi3.Parameters, operationId string) error {
	l := oa.logger()
	level := utils.LogLevel(behaviour)
	if paramsType == nilType {
		return nil
	}
	validator := schema_validator.NewTypeSchemaValidator(l, level, nilType, openapi3.Schema{})
	for _, field := range reflect.VisibleFields(paramsType) {
		name := field.Tag.Get(tag)
		if name == ignoreFieldTagValue {
			continue
		}
		if name == "" {
			name = field.Name
		}
		specParameter := specParameters.GetByInAndName(in, name)
		if specParameter == nil {
			l.Logf(level, paramDefinedByHandlerButMissingInSpec(in, name, paramsType, operationId))
			continue
		}
		// TODO: schema validator check object schemas with json keys
		l.LogIfNotNil(level, validator.WithType(field.Type).WithSchema(*specParameter.Schema.Value).Validate())
	}
	return l.MustHaveNoErrors()
}

// validateResponseTypes check that all responses declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all possible responses defined in the spec, but it can not declare responses which are not defined.
func validateResponseTypes(oa openapi, behaviour Behaviour, handler handler, specOperation *openapi3.Operation, operationId string) error {
	l := oa.logger()
	level := utils.LogLevel(behaviour)
	validator := schema_validator.NewTypeSchemaValidator(l, level, nilType, openapi3.Schema{})
	for status, response := range handler.responses {
		specResponse := specOperation.Responses.Get(status)
		if specResponse == nil {
			l.Logf(level, handlerDefinesResponseThatIsMissingInTheSpec(status, operationId))
			continue
		}
		responseValidator := validator.WithType(response.responseType)
		for _, mediaType := range specResponse.Value.Content {
			// TODO: allow different media types to fine tune behaviour of validation (for example use non json tags during struct validation for non json mime type)
			l.LogIfNotNil(level, responseValidator.WithSchema(*mediaType.Schema.Value).Validate())
		}
	}
	return l.MustHaveNoErrors()
}
