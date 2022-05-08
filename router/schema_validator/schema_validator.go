package schema_validator

import (
	"encoding"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router/utils"
	"reflect"
)

type Options struct {
}

var textMarshallerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

// schema types allowed by OpenAPI specification.
const (
	objectSchemaType  = "object"
	arraySchemaType   = "array"
	stringSchemaType  = "string"
	booleanSchemaType = "boolean"
	numberSchemaType  = "number"
	integerSchemaType = "integer"
)

// TypeSchemaValidator helps validate reflect.Type and openapi3.Schema compatibility using the validation Options.
type TypeSchemaValidator interface {
	// WithOptions immutably returns a new TypeSchemaValidator with the specified validation options.
	WithOptions(Options) TypeSchemaValidator
	// WithType immutably returns a new TypeSchemaValidator with the specified reflect.Type to validate.
	WithType(reflect.Type) TypeSchemaValidator
	// WithSchema immutably returns a new TypeSchemaValidator with the specified openapi3.Schema to validate.
	WithSchema(openapi3.Schema) TypeSchemaValidator
	// Validate reflect.Type and the openapi3.Schema compatibility using the validation Options.
	// Returns utils.MultiError with all compatability errors found or nil if compatible.
	Validate() utils.MultiError

	validateSchemaAllOf() utils.MultiError
	validateSchemaOneOf() utils.MultiError
	validateSchemaAnyOf() utils.MultiError
	validateSchemaNot() utils.MultiError
	validateObjectSchema() utils.MultiError
	validateArraySchema() utils.MultiError
	validateStringSchema() utils.MultiError
	validateBooleanSchema() utils.MultiError
	validateIntegerSchema() utils.MultiError
	validateNumberSchema() utils.MultiError
}

// NewTypeSchemaValidator returns a new TypeSchemaValidator that helps validate reflect.Type and openapi3.Schema compatibility using the validation Options.
func NewTypeSchemaValidator(goType reflect.Type, schema openapi3.Schema, options Options) TypeSchemaValidator {
	return typeSchemaValidatorContext{
		goType:  goType,
		schema:  schema,
		options: options,
	}
}

// typeSchemaValidatorContext an internal struct that implementation TypeSchemaValidator
type typeSchemaValidatorContext struct {
	goType  reflect.Type
	schema  openapi3.Schema
	options Options
}

func (c typeSchemaValidatorContext) WithOptions(options Options) TypeSchemaValidator {
	c.options = options
	return c
}
func (c typeSchemaValidatorContext) WithType(goType reflect.Type) TypeSchemaValidator {
	c.goType = goType
	return c
}
func (c typeSchemaValidatorContext) WithSchema(schema openapi3.Schema) TypeSchemaValidator {
	c.schema = schema
	return c
}
func (c typeSchemaValidatorContext) Validate() utils.MultiError {
	errs := utils.NewErrorsCollector()

	// Test global schema validation properties
	errs.AddErrorsIfNotNil(c.validateSchemaAllOf())
	errs.AddErrorsIfNotNil(c.validateSchemaOneOf())
	errs.AddErrorsIfNotNil(c.validateSchemaAnyOf())
	errs.AddErrorsIfNotNil(c.validateSchemaNot())

	// Test specific schema types validations
	switch c.schema.Type {
	case objectSchemaType:
		errs.AddErrorsIfNotNil(c.validateObjectSchema())
	case arraySchemaType:
		errs.AddErrorsIfNotNil(c.validateArraySchema())
	case stringSchemaType:
		errs.AddErrorsIfNotNil(c.validateStringSchema())
	case booleanSchemaType:
		errs.AddErrorsIfNotNil(c.validateBooleanSchema())
	case numberSchemaType:
		errs.AddErrorsIfNotNil(c.validateNumberSchema())
	case integerSchemaType:
		errs.AddErrorsIfNotNil(c.validateIntegerSchema())
	}
	return errs.ErrorOrNil()
}
