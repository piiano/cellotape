package schema_validator

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/utils"
)

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
	// WithLogger immutably returns a new TypeSchemaValidator with the specified utils.Logger.
	WithLogger(utils.Logger) TypeSchemaValidator
	// WithType immutably returns a new TypeSchemaValidator with the specified reflect.Type to validate.
	WithType(reflect.Type) TypeSchemaValidator
	// WithSchema immutably returns a new TypeSchemaValidator with the specified openapi3.Schema to validate.
	WithSchema(openapi3.Schema) TypeSchemaValidator
	// WithSchemaAndType immutably returns a new TypeSchemaValidator with the specified openapi3.Schema and reflect.Type to validate.
	WithSchemaAndType(openapi3.Schema, reflect.Type) TypeSchemaValidator
	// Validate reflect.Type and the openapi3.Schema compatibility using the validation Options.
	// Returns error with all compatability errors found or nil if compatible.
	Validate() error

	validateSchemaAllOf() error
	validateSchemaOneOf() error
	validateSchemaAnyOf() error
	validateSchemaNot() error
	validateObjectSchema() error
	validateArraySchema() error
	validateStringSchema() error
	validateBooleanSchema() error
	validateIntegerSchema() error
	validateNumberSchema() error

	newLogger() utils.Logger
	logLevel() utils.LogLevel
}

// NewEmptyTypeSchemaValidator returns a new TypeSchemaValidator that have no reflect.Type or openapi3.Schema configured yet.
func NewEmptyTypeSchemaValidator(logger utils.Logger) TypeSchemaValidator {
	return typeSchemaValidatorContext{
		logger: logger,
		level:  utils.Error,
	}
}

// NewTypeSchemaValidator returns a new TypeSchemaValidator that helps validate reflect.Type and openapi3.Schema compatibility using the validation Options.
func NewTypeSchemaValidator(logger utils.Logger, level utils.LogLevel, goType reflect.Type, schema openapi3.Schema) TypeSchemaValidator {
	return typeSchemaValidatorContext{
		logger: logger,
		level:  level,
		schema: schema,
		goType: goType,
	}
}

func (c typeSchemaValidatorContext) newLogger() utils.Logger {
	return c.logger.NewCounter()
}
func (c typeSchemaValidatorContext) logLevel() utils.LogLevel {
	return c.level
}

// typeSchemaValidatorContext an internal struct that implementation TypeSchemaValidator
type typeSchemaValidatorContext struct {
	logger utils.Logger
	level  utils.LogLevel
	schema openapi3.Schema
	goType reflect.Type
}

func (c typeSchemaValidatorContext) WithLogger(logger utils.Logger) TypeSchemaValidator {
	c.logger = logger
	return c
}
func (c typeSchemaValidatorContext) WithType(goType reflect.Type) TypeSchemaValidator {
	c.goType = goType
	c.logger = c.newLogger()
	return c
}
func (c typeSchemaValidatorContext) WithSchema(schema openapi3.Schema) TypeSchemaValidator {
	c.schema = schema
	c.logger = c.newLogger()
	return c
}
func (c typeSchemaValidatorContext) WithSchemaAndType(schema openapi3.Schema, goType reflect.Type) TypeSchemaValidator {
	c.schema = schema
	c.goType = goType
	c.logger = c.newLogger()
	return c
}

func (c typeSchemaValidatorContext) Validate() error {
	if isEmptyInterface(c.goType) {
		return nil
	}
	if c.goType.Kind() == reflect.Pointer {
		return c.WithType(c.goType.Elem()).Validate()
	}
	// Test global schema validation properties
	c.logger.ErrorIfNotNil(c.validateSchemaAllOf())
	c.logger.ErrorIfNotNil(c.validateSchemaOneOf())
	c.logger.ErrorIfNotNil(c.validateSchemaAnyOf())
	c.logger.ErrorIfNotNil(c.validateSchemaNot())

	// Test specific schema types validations
	switch c.schema.Type {
	case objectSchemaType:
		c.logger.ErrorIfNotNil(c.validateObjectSchema())
	case arraySchemaType:
		c.logger.ErrorIfNotNil(c.validateArraySchema())
	case stringSchemaType:
		c.logger.ErrorIfNotNil(c.validateStringSchema())
	case booleanSchemaType:
		c.logger.ErrorIfNotNil(c.validateBooleanSchema())
	case numberSchemaType:
		c.logger.ErrorIfNotNil(c.validateNumberSchema())
	case integerSchemaType:
		c.logger.ErrorIfNotNil(c.validateIntegerSchema())
	}
	return c.logger.MustHaveNoErrors()
}

func isEmptyInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface && t.NumMethod() == 0
}
