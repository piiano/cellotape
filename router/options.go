package router

import "github.com/piiano/restcontroller/schema_validator"

type Options struct {
	InitializationSchemaValidation schema_validator.Options
	//	TODO: Add support for options
	//	// Return error in OpenAPI.Validate if there are missing implementation for components of the OpenAPISpec.
	//	FailOnMissing struct {
	//		// Return error for operations on the OpenAPISpec with no matching OperationHandler implementation.
	//		Operations bool
	//		// Return error for responses on the OpenAPISpec with no matching HttpResponse implementation.
	//		Responses bool
	//	}
	//	// Allow tuning schema validation during initialization.
	//	InitializationSchemaValidation struct {
	//		// When true, enable reflection validation of schema compatability with go types during initialization.
	//		Enable bool
	//		// Strictly disallow use of empty interface (any) type.
	//		AllowAny bool
	//		// Strictly require that unsigned integers will have a schema min validation.
	//		UintMustHaveZeroMinValue bool
	//		// Strictly require that schema max validation will be compatible with the type capacity.
	//		ValidateMaxCompatibleWithTypeSize bool
	//	}
	//	// Enforce schema validation at runtime using the OpenAPISpec schema.
	//	RuntimeSchemaValidation struct {
	//		// Enforce schema validation at runtime for RequestBody.
	//		RequestBody bool
	//		// Enforce schema validation at runtime for ResponseBody.
	//		ResponseBody bool
	//	}
}
