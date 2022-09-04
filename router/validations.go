package router

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/schema_validator"
	"github.com/piiano/cellotape/router/utils"
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
	declaredOperation := utils.NewSet[string]()
	excludeOperations := utils.NewSet(oa.options.ExcludeOperations...)
	l.ErrorIfNotNil(validateContentTypes(*oa, excludeOperations))
	for _, flatOp := range flatOperations {
		if excludeOperations.Has(flatOp.id) {
			l.Errorf(anExcludedOperationIsImplemented(flatOp.id))
		}
		if !declaredOperation.Add(flatOp.id) {
			// multiple handlers for the same operation is always an error
			l.Errorf(multipleHandlersFoundForOperationId(flatOp.id))
		}
		l.ErrorIfNotNil(validateOperation(*oa, flatOp))
	}
	l.ErrorIfNotNil(validateMustHandleAllOperations(oa, declaredOperation, excludeOperations))
	return l.MustHaveNoErrorsf(failedValidatingTheRouterWithTheSpec(l.Warnings(), l.Errors()))
}

// validateMustHandleAllOperations checks that all operations defined in the spec have an implementation on the router.
func validateMustHandleAllOperations(oa *openapi, declaredOperations utils.Set[string], excludeOperations utils.Set[string]) error {
	l := oa.logger()
	for _, pathItem := range oa.spec.Paths {
		for _, specOp := range pathItem.Operations() {
			if excludeOperations.Has(specOp.OperationID) {
				continue
			}
			if !declaredOperations.Has(specOp.OperationID) {
				l.Logf(utils.LogLevel(oa.options.MustHandleAllOperations), missingHandlerForOperationId(specOp.OperationID))
			}
		}
	}
	return l.MustHaveNoErrorsf(notImplementedSpecOperations(l.Errors()))
}

// validateContentTypes checks that all content types defined in the spec for request or responses have an implementation on the router.
func validateContentTypes(oa openapi, excludeOperations utils.Set[string]) error {
	log := oa.logger()
	level := utils.LogLevel(oa.options.HandleAllContentTypes)
	specContentTypes := oa.spec.findSpecContentTypes(excludeOperations)
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
	handlersChain := append(operation.handlers, operation.handler)
	for _, chainHandler := range handlersChain {
		l.AppendCounters(validateRequestBodyType(oa, options.ValidateRequestBody, chainHandler, specOp.RequestBody, operation.id))
		l.AppendCounters(validatePathParamsType(oa, options.ValidatePathParams, chainHandler, specOp.Parameters, operation.id))
		l.AppendCounters(validateQueryParamsType(oa, options.ValidateQueryParams, chainHandler, specOp.Parameters, operation.id))
		l.AppendCounters(validateResponseTypes(oa, options.ValidateResponses, chainHandler, specOp.Operation, operation.id))
	}
	l.AppendCounters(validateHandleAllPathParams(oa, options.HandleAllPathParams, operation, specOp))
	l.AppendCounters(validateHandleAllQueryParams(oa, options.HandleAllQueryParams, operation, specOp))
	l.AppendCounters(validateHandleAllResponses(oa, options.HandleAllOperationResponses, operation, specOp))
	return l.MustHaveNoErrorsf("operation %q has incompatibility with the spec (%d errors, %d warnings)", operation.id, l.Errors(), l.Warnings())
}

// validateHandleAllPathParams checks that every path param defined in the operation is handled at least once in the handlers chain
func validateHandleAllPathParams(oa openapi, behaviour Behaviour, operation operation, specOp SpecOperation) utils.LogCounters {
	handlers := append(operation.handlers, operation.handler)
	declaredParams := utils.NewSet[string](utils.ConcatSlices[string](utils.Map(handlers, func(h handler) []string {
		return utils.Keys(structKeys(h.request.pathParams, pathParamFieldTag))
	})...)...)
	return validateHandleAllParams(oa, behaviour, operation, specOp, pathParamInValue, declaredParams)
}

// validateHandleAllQueryParams checks that every query param defined in the operation is handled at least once in the handlers chain
func validateHandleAllQueryParams(oa openapi, behaviour Behaviour, operation operation, specOp SpecOperation) utils.LogCounters {
	handlers := append(operation.handlers, operation.handler)
	declaredParams := utils.NewSet[string](utils.ConcatSlices[string](utils.Map(handlers, func(h handler) []string {
		return utils.Keys(structKeys(h.request.queryParams, queryParamFieldTag))
	})...)...)
	return validateHandleAllParams(oa, behaviour, operation, specOp, queryParamInValue, declaredParams)
}

// validateHandleAllParams checks that every parameter defined in the operation is handled at least once in the handlers chain
func validateHandleAllParams(oa openapi, behaviour Behaviour, operation operation, specOp SpecOperation, in string, declaredParams utils.Set[string]) utils.LogCounters {
	l := oa.logger()
	level := utils.LogLevel(behaviour)
	for _, specParam := range specOp.Parameters {
		if specParam.Value.In != in {
			continue
		}
		name := specParam.Value.Name
		if !declaredParams.Has(name) {
			l.Logf(level, paramMissingImplementationInChain(in, name, operation.id))
		}
	}
	return l.Counters()
}

// validateHandleAllResponses checks that every response defined in the spec is handled at least once in the handlers chain
func validateHandleAllResponses(oa openapi, behaviour Behaviour, operation operation, specOp SpecOperation) utils.LogCounters {
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
			l.Logf(level, unimplementedResponseForOperation(status, operation.id))
		}
	}
	return l.Counters()
}

// validateRequestBodyType check that a request body type declared on a handler is declared on the spec with a compatible schema.
// a handler does not have to declare and handle the request body defined in the spec, but it can not declare request body which is not defined or incompatible.
func validateRequestBodyType(oa openapi, behaviour Behaviour, handler handler, specBody *openapi3.RequestBodyRef, operationID string) utils.LogCounters {
	l := oa.logger()
	level := utils.LogLevel(behaviour)
	bodyType := handler.request.requestBody
	if bodyType == nilType {
		return utils.LogCounters{}
	}
	validator := schema_validator.NewTypeSchemaValidator(l, level, bodyType, openapi3.Schema{})
	if specBody == nil {
		l.Logf(level, handlerDefinesRequestBodyWhenNoRequestBodyInSpec(operationID))
		return l.Counters()
	}
	for _, mediaType := range specBody.Value.Content {
		// TODO: allow different media types to fine tune behaviour of validation (for example use non json tags during struct validation for non json mime type)
		if err := validator.WithSchema(*mediaType.Schema.Value).Validate(); err != nil {
			l.Logf(level, incompatibleRequestBodyType(operationID, bodyType))
		}
	}
	return l.Counters()
}

// validatePathParamsType check that all pathParamInValue params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all pathParamInValue parameters defined in the spec, but it can not declare parameters which are not defined.
func validatePathParamsType(oa openapi, behaviour Behaviour, handler handler, specParameters openapi3.Parameters, operationId string) utils.LogCounters {
	return validateParamsType(oa, behaviour, pathParamInValue, pathParamFieldTag, handler.request.pathParams, specParameters, operationId)
}

// validatePathParamsType check that all queryParamInValue params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all queryParamInValue parameters defined in the spec, but it can not declare parameters which are not defined.
func validateQueryParamsType(oa openapi, behaviour Behaviour, handler handler, specParameters openapi3.Parameters, operationId string) utils.LogCounters {
	return validateParamsType(oa, behaviour, queryParamInValue, queryParamFieldTag, handler.request.queryParams, specParameters, operationId)
}

// validateParamsType check that all params declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all parameters defined in the spec, but it can not declare parameters which are not defined.
func validateParamsType(oa openapi, behaviour Behaviour, in string, tag string, paramsType reflect.Type, specParameters openapi3.Parameters, operationId string) utils.LogCounters {
	l := oa.logger()
	level := utils.LogLevel(behaviour)
	if paramsType == nilType {
		return utils.LogCounters{}
	}
	validator := schema_validator.NewTypeSchemaValidator(l, level, nilType, openapi3.Schema{})
	for name, field := range structKeys(paramsType, tag) {
		specParameter := specParameters.GetByInAndName(in, name)
		if specParameter == nil {
			l.Logf(level, paramDefinedByHandlerButMissingInSpec(in, name, paramsType, operationId))
			continue
		}
		// TODO: schema validator check object schemas with json keys
		if err := validator.WithType(field.Type).WithSchema(*specParameter.Schema.Value).Validate(); err != nil {
			l.Logf(level, incompatibleParamType(operationId, in, name, field.Name, field.Type))
		}
	}
	return l.Counters()
}

func structKeys(structType reflect.Type, tag string) map[string]reflect.StructField {
	if structType == nil || structType == nilType {
		return map[string]reflect.StructField{}
	}
	return utils.FromEntries(utils.ConcatSlices(utils.Map(utils.Filter(reflect.VisibleFields(structType), func(field reflect.StructField) bool {
		return !field.Anonymous && field.IsExported() && field.Tag.Get(tag) != ignoreFieldTagValue
	}), func(field reflect.StructField) utils.Entry[string, reflect.StructField] {
		name := field.Tag.Get(tag)
		name, _, _ = strings.Cut(name, ",")
		if name == "" {
			name = field.Name
		}
		return utils.Entry[string, reflect.StructField]{Key: name, Value: field}
	}), utils.ConcatSlices(utils.Map(utils.Filter(reflect.VisibleFields(structType), func(field reflect.StructField) bool {
		return field.Anonymous && field.IsExported() && field.Tag.Get(tag) != ignoreFieldTagValue
	}), func(field reflect.StructField) []utils.Entry[string, reflect.StructField] {
		return utils.Entries(structKeys(field.Type, tag))
	})...)))
}

// validateResponseTypes check that all responses declared on a handler are available on the spec with a compatible schema.
// a handler does not have to declare and handle all possible responses defined in the spec, but it can not declare responses which are not defined.
func validateResponseTypes(oa openapi, behaviour Behaviour, handler handler, specOperation *openapi3.Operation, operationId string) utils.LogCounters {
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
			if err := responseValidator.WithSchema(*mediaType.Schema.Value).Validate(); err != nil {
				l.Logf(level, incompatibleResponseType(operationId, status, response.responseType))
			}
		}
	}
	return l.Counters()
}
