package router

import "github.com/piiano/restcontroller/schema_validator"

type Options struct {
	// TODO: Add support for options
	// Fine tune schema validation during initialization.
	InitializationSchemaValidation schema_validator.Options

	// provide a default recover from panic that return status 500
	RecoverOnPanic bool
}
