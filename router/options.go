package router

import (
	"fmt"
	"io"
	"os"

	"github.com/piiano/cellotape/router/utils"
)

type LogLevel = utils.LogLevel

const (
	LogLevelError = utils.Error
	LogLevelWarn  = utils.Warn
	LogLevelInfo  = utils.Info
	LogLevelOff   = utils.Off
)

// Behaviour defines a possible behaviour for a validation error.
// Possible values are PropagateError, PrintWarning and Ignore.
//
// 		- PropagateError - Cause call to OpenAPIRouter.AsHandler to return an error if the validation fails and failure reason is logged with error level.
// 		- PrintWarning - if the validation fails the failure reason will get logged with warning level. The error won't propagate to OpenAPIRouter.AsHandler return values
// 		- Ignore - do nothing if the validation fails. don't print anything to the logs and don't propagate the error to OpenAPIRouter.AsHandler return values
//
// By default, Behaviour are initialized to PropagateError as it is essentially Behaviour zero value.
type Behaviour utils.LogLevel

const (
	// PropagateError propagate validation errors to OpenAPIRouter.AsHandler return values and print them to the log.
	// This is also the initial Behaviour value if not set to anything else.
	PropagateError = Behaviour(utils.Error)
	// PrintWarning print validation errors as warning to the log without propagating the error to OpenAPIRouter.AsHandler return values
	PrintWarning = Behaviour(utils.Warn)
	// Ignore do nothing if the validation fails.
	// A validation failure won't get printed to the logs and an error won't get propagate to OpenAPIRouter.AsHandler return values
	Ignore = Behaviour(utils.Off)
)

// MarshalText implements encoding.TextMarshaler.
func (b Behaviour) MarshalText() ([]byte, error) {
	switch b {
	case Ignore:
		return []byte("ignore"), nil
	case PrintWarning:
		return []byte("print-warning"), nil
	case PropagateError:
		return []byte("propagate-error"), nil
	}
	return nil, fmt.Errorf("%d is an invalid behaviour value", b)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Behaviour) UnmarshalText(data []byte) error {
	behaviourString := string(data)
	switch behaviourString {
	case "ignore":
		*b = Ignore
	case "print-warning":
		*b = PrintWarning
	case "propagate-error":
		*b = PropagateError
	default:
		return fmt.Errorf("%s is an invalid behaviour value", behaviourString)
	}
	return nil
}

// Options defines the behaviour of the OpenAPI router
type Options struct {

	// When RecoverOnPanic is set to true the handlers chain provide a default recover behaviour that return status 500.
	// This feature is available and turned on by default to prevent the server from crashing when an internal panic
	// occurs in the handler chain.
	// This is a fundamental requirement of a web server to support availability.
	// While being fundamental, the default behaviour can't conform to your spec out of the box.
	// To be able to better control the behaviour during panic and the structure of the response returned to be
	// compatible with the spec it's recommended to implement it as a handler and in this case the default recover
	// behaviour can be turned off.
	RecoverOnPanic bool `json:"recoverOnPanic,omitempty"`

	// LogLevel defines what log levels should be printed to LogOutput.
	// By default, LogLevel is set to utils.Info to print all info to the log.
	// The router prints to the log only during initialization to show validation errors, warnings and useful info.
	// No printing is done after initialization.
	LogLevel LogLevel `json:"logLevel,omitempty"`

	// LogOutput defines where to write the outputs too.
	// By default, it is set to write to os.Stderr.
	// You can provide your own writer to be able to read the validation errors programmatically or to write them to
	// different destination.
	// The router prints to the log only during initialization to show validation errors, warnings and useful info.
	// No printing is done after initialization.
	LogOutput io.Writer `json:"-"`

	// OperationValidations allow defining validation for specific operations using a map of operation id to an
	// operationValidationOptions structure.
	// This option is used only to override the default operation validations options defined ny the
	// DefaultOperationValidation option.
	OperationValidations map[string]OperationValidationOptions `json:"operationValidations,omitempty"`

	// DefaultOperationValidation defines the default validations run for every operation.
	// By default, if not set the default behaviour is strict and consider any validation error to fail the entire
	// router validations.
	DefaultOperationValidation OperationValidationOptions `json:"defaultOperationValidation,omitempty"`

	// The router validates that every operation defined in the spec has an implementation.
	// MustHandleAllOperations defines the behaviour when this validation fails.
	// By default, it is set to utils.Error to propagate the error.
	// If your implementation is still in progress, and you would like to ignore unimplemented operations you can set
	// this option to utils.Warn or turn it off with utils.Off.
	MustHandleAllOperations Behaviour `json:"mustHandleAllOperations,omitempty"`

	// The router validates that every content type define in the spec has an implementation.
	// HandleAllContentTypes defines the behaviour when this validation fails.
	// By default, it is set to utils.Error to propagate the error.
	// Content types defined in the spec in operation request body and responses.
	// If your spec defines some content types that are not supported by the default router content types you can add
	// your own implementation.
	// Content type implementation is basically an implementation of ContentType interface add with
	// OpenAPIRouter.WithContentType to define serialization and deserialization behaviour.
	HandleAllContentTypes Behaviour `json:"handleAllContentTypes,omitempty"`

	// ExcludeOperations defined an array of operations that are defined in the spec but are excluded from the implementation.
	// One use for this option is when your spec defines the entire API of your app but the implementation is spread to multiple microservices.
	// With this option you can define list of operations that are to be implemented by other microservices.
	ExcludeOperations []string
}

// OperationValidationOptions defines options to control operation validations
type OperationValidationOptions struct {
	// ValidatePathParams determines validation of operation request body.
	ValidateRequestBody Behaviour `json:"validateRequestBody,omitempty"`

	// ValidatePathParams determines validation of operation path params.
	ValidatePathParams Behaviour `json:"validatePathParams,omitempty"`

	// HandleAllPathParams describes the behaviour when not every path params defined in the spec is handled at least once in the handlers chain
	HandleAllPathParams Behaviour `json:"handleAllPathParams,omitempty"`

	// ValidatePathParams determines validation of operation query params.
	ValidateQueryParams Behaviour `json:"validateQueryParams,omitempty"`

	// HandleAllQueryParams describes the behaviour when not every query params defined in the spec is handled at least once in the handlers chain
	HandleAllQueryParams Behaviour `json:"handleAllQueryParams,omitempty"`

	// ValidatePathParams determines validation of operation responses.
	ValidateResponses Behaviour `json:"validateResponses,omitempty"`

	// HandleAllOperationResponses describes the behaviour when not every response defined in the spec is handled at least once in the handlers chain
	HandleAllOperationResponses Behaviour `json:"handleAllOperationResponses,omitempty"`
}

// SchemaValidationOptions defines options to control schema validations
type SchemaValidationOptions struct {
	// NoEmptyInterface defines the behaviour when a schema is validated with an empty interface (any) type.
	NoEmptyInterface Behaviour `json:"noEmptyInterface,omitempty"`

	// NoStringAnyMapForObjectsSchema defines the behaviour when an object schema is validated with a string to empty interface map type (map[string]any).
	NoStringAnyMapForObjectsSchema Behaviour `json:"noStringAnyMapForObjectsSchema,omitempty"`
}

// DefaultOptions returns the default OpenAPI Router Options.
// For loading different options it is recommended to start from the default Options as the baseline and modify the
// specific options needed.
// To init the router with custom Options you can use NewOpenAPIRouterWithOptions.
func DefaultOptions() Options {
	return Options{
		RecoverOnPanic: true,
		LogLevel:       LogLevelInfo,
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

func (o Options) operationValidationOptions(id string) OperationValidationOptions {
	options, ok := o.OperationValidations[id]
	if ok {
		return options
	}
	return o.DefaultOperationValidation
}
