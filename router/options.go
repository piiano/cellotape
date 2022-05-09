package router

import (
	"github.com/piiano/restcontroller/router/utils"
	"io"
	"os"
)

// Options defines the behaviour of the OpenAPI router
type Options struct {
	// TODO: Add support fine tuning schema validation during initialization.

	// When RecoverOnPanic is set to true the handlers chain provide a default recover behaviour that return status 500.
	// This feature is available and turned on by default to prevent the server from crashing when an internal panic
	// occurs in the handler chain.
	// This is a fundamental requirement of a web server to support availability.
	// While being fundamental, the default behaviour can't conform to your spec out of the box.
	// To be able to better control the behaviour during panic and the structure of the response returned to be
	// compatible with the spec it's recommended to implement it as a handler and in this case the default recover
	// behaviour can be turned off.
	RecoverOnPanic bool

	// LogLevel defines what log levels should be printed to LogOutput.
	// By default, LogLevel is set to utils.Info to print all info to the log.
	// The router prints to the log only during initialization to show validation errors, warnings and useful info.
	// No printing is done after initialization.
	LogLevel utils.LogLevel

	// LogOutput defines where to write the outputs too.
	// By default, it is set to write to os.Stderr.
	// You can provide your own writer to be able to read the validation errors programmatically or to write them to
	// different destination.
	// The router prints to the log only during initialization to show validation errors, warnings and useful info.
	// No printing is done after initialization.
	LogOutput io.Writer

	// OperationValidations allow defining validation for specific operations using a map of operation id to an
	// OperationValidationOptions structure.
	// This option is used only to override the default operation validations options defined ny the
	// DefaultOperationValidation option.
	OperationValidations map[string]OperationValidationOptions

	// DefaultOperationValidation defines the default validations run for every operation.
	// By default, if not set the default behaviour is strict and consider any validation error to fail the entire
	// router validations.
	DefaultOperationValidation OperationValidationOptions

	// The router validates that every operation defined in the spec has an implementation.
	// MustHandleAllOperations defines the behaviour when this validation fails.
	// By default, it is set to utils.Error to propagate the error.
	// If your implementation is still in progress, and you would like to ignore unimplemented operations you can set
	// this option to utils.Warn or turn it off with utils.Off.
	MustHandleAllOperations Behaviour

	// The router validates that every content type define in the spec has an implementation.
	// HandleAllContentTypes defines the behaviour when this validation fails.
	// By default, it is set to utils.Error to propagate the error.
	// Content types defined in the spec in operation request body and responses.
	// If your spec defines some content types that are not supported by the default router content types you can add
	// your own implementation.
	// Content type implementation is basically an implementation of ContentType interface add with
	// OpenAPIRouter.WithContentType to define serialization and deserialization behaviour.
	HandleAllContentTypes Behaviour
}

type Behaviour utils.LogLevel

const (
	// PropagateError propagate validation errors to the next validation layer as a returned error.
	// This is also the initial Behaviour value if not set to anything else.
	PropagateError = Behaviour(utils.Error)
	// PrintWarning print validation errors as warning to the log
	PrintWarning = Behaviour(utils.Warn)
	// Off ignore validation errors
	Off = Behaviour(utils.Off)
)

// OperationValidationOptions defines options to control operation validations
type OperationValidationOptions struct {
	// ValidatePathParams determines validation of operation request body.
	ValidateRequestBody Behaviour

	// ValidatePathParams determines validation of operation path params.
	ValidatePathParams Behaviour

	// ValidatePathParams determines validation of operation query params.
	ValidateQueryParams Behaviour

	// ValidatePathParams determines validation of operation responses.
	ValidateResponses Behaviour

	// When HandleAllOperationResponses is set to true there is a check that every response defined in the spec is handled at least once in the handlers chain
	HandleAllOperationResponses Behaviour
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
		LogLevel:       utils.Info,
		LogOutput:      os.Stderr,
		DefaultOperationValidation: OperationValidationOptions{
			ValidateRequestBody:         PropagateError,
			ValidatePathParams:          PropagateError,
			ValidateQueryParams:         PropagateError,
			ValidateResponses:           PropagateError,
			HandleAllOperationResponses: PropagateError,
		},
		MustHandleAllOperations: PropagateError,
		HandleAllContentTypes:   PropagateError,
	}
}
