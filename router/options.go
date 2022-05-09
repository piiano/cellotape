package router

import (
	"github.com/piiano/restcontroller/router/schema_validator"
	"github.com/piiano/restcontroller/router/utils"
)

type Options struct {
	// TODO: Add support for options
	// Fine tune schema validation during initialization.
	SchemaValidation schema_validator.Options

	// When RecoverOnPanic is set to true the handlers chain provide a default recover behaviour that return status 500.
	RecoverOnPanic bool

	// DefaultOperationValidation defines the default validations run for every operation
	DefaultOperationValidation OperationValidationOptions

	// OperationValidations allow overriding the validations defined ny DefaultOperationValidation for specific operations using their operation id
	OperationValidations map[string]OperationValidationOptions

	// HandleAllOperations is set to true there is a check that every operation defined in the spec has an implementation in the router
	HandleAllOperations utils.LogLevel

	// HandleAllContentTypes is set to true there is a check that every operation defined in the spec has an implementation in the router
	HandleAllContentTypes utils.LogLevel
}

const (
	ReturnError   = utils.Error
	PrintWarnning = utils.Warn
	Ignore        = utils.Ignore
)

type OperationValidationOptions struct {
	// ValidatePathParams determines validation of operation request body.
	ValidateRequestBody utils.LogLevel

	// ValidatePathParams determines validation of operation path params.
	ValidatePathParams utils.LogLevel

	// ValidatePathParams determines validation of operation query params.
	ValidateQueryParams utils.LogLevel

	// ValidatePathParams determines validation of operation responses.
	ValidateResponses utils.LogLevel

	// When HandleAllOperationResponses is set to true there is a check that every response defined in the spec is handled at least once in the handlers chain
	HandleAllOperationResponses utils.LogLevel
}

func (o Options) OperationValidationOptions(id string) OperationValidationOptions {
	options, ok := o.OperationValidations[id]
	if ok {
		return options
	}
	return o.DefaultOperationValidation
}

func DefaultOptions() Options {
	return Options{
		RecoverOnPanic: true,
	}
}
