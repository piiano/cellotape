package schema_validator

import (
	"encoding"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/utils"
	"reflect"
)

type Options struct {
}

var textMarshallerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

type SchemaType string

// schema types allowed by OpenAPI specification.
const (
	objectSchemaType  SchemaType = "object"
	arraySchemaType   SchemaType = "array"
	stringSchemaType  SchemaType = "string"
	booleanSchemaType SchemaType = "boolean"
	numberSchemaType  SchemaType = "number"
	integerSchemaType SchemaType = "integer"
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
	errs.AddIfNotNil(c.validateSchemaAllOf())
	errs.AddIfNotNil(c.validateSchemaOneOf())
	errs.AddIfNotNil(c.validateSchemaAnyOf())
	errs.AddIfNotNil(c.validateSchemaNot())

	// Test specific schema types validations
	switch SchemaType(c.schema.Type) {
	case objectSchemaType:
		errs.AddIfNotNil(c.validateObjectSchema())
	case arraySchemaType:
		errs.AddIfNotNil(c.validateArraySchema())
	case stringSchemaType:
		errs.AddIfNotNil(c.validateStringSchema())
	case booleanSchemaType:
		errs.AddIfNotNil(c.validateBooleanSchema())
	case numberSchemaType:
		errs.AddIfNotNil(c.validateNumberSchema())
	case integerSchemaType:
		errs.AddIfNotNil(c.validateIntegerSchema())
	}
	return errs.ErrorOrNil()
}
