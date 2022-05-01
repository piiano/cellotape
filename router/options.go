package router

import "github.com/piiano/restcontroller/schema_validator"

type Options struct {
	// TODO: Add support for options
	// Allow tuning schema validation during initialization.
	InitializationSchemaValidation schema_validator.Options
}
