package schema_validator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/utils"
)

var ErrSchemaIncompatibleWithType = errors.New("schema is incompatible with type")

// TypeSchemaValidator helps validate reflect.Type and openapi3.Schema compatibility using the validation Options.
type TypeSchemaValidator interface {
	// WithType immutably returns a new TypeSchemaValidator with the specified reflect.Type to validate.
	WithType(reflect.Type) TypeSchemaValidator
	// WithSchema immutably returns a new TypeSchemaValidator with the specified openapi3.Schema to validate.
	WithSchema(openapi3.Schema) TypeSchemaValidator
	// WithSchemaAndType immutably returns a new TypeSchemaValidator with the specified openapi3.Schema and reflect.Type to validate.
	WithSchemaAndType(openapi3.Schema, reflect.Type) TypeSchemaValidator
	// Validate reflect.Type and the openapi3.Schema compatibility using the validation Options.
	// Returns error with all compatability errors found or nil if compatible.
	Validate() error

	Errors() []string

	matchAllSchemaValidator(string, openapi3.SchemaRefs)
	validateSchemaAllOf()
	validateSchemaNot()
	validateObjectSchema()
	validateArraySchema()
	validateStringSchema()
	validateBooleanSchema()
	validateIntegerSchema()
	validateNumberSchema()
}

// NewEmptyTypeSchemaValidator returns a new TypeSchemaValidator that have no reflect.Type or openapi3.Schema configured yet.
func NewEmptyTypeSchemaValidator() TypeSchemaValidator {
	return typeSchemaValidatorContext{
		errors: new([]string),
	}
}

// NewTypeSchemaValidator returns a new TypeSchemaValidator that helps validate reflect.Type and openapi3.Schema compatibility using the validation Options.
func NewTypeSchemaValidator(goType reflect.Type, schema openapi3.Schema) TypeSchemaValidator {
	return typeSchemaValidatorContext{
		errors: new([]string),
		schema: schema,
		goType: goType,
	}
}

// typeSchemaValidatorContext an internal struct that implementation TypeSchemaValidator
type typeSchemaValidatorContext struct {
	errors *[]string
	schema openapi3.Schema
	goType reflect.Type
}

func (c typeSchemaValidatorContext) err(format string, args ...any) {
	*c.errors = append(*c.errors, fmt.Sprintf(format, args...))
}

func (c typeSchemaValidatorContext) WithType(goType reflect.Type) TypeSchemaValidator {
	c.goType = goType
	return c
}
func (c typeSchemaValidatorContext) WithSchema(schema openapi3.Schema) TypeSchemaValidator {
	c.schema = schema
	return c
}
func (c typeSchemaValidatorContext) WithSchemaAndType(schema openapi3.Schema, goType reflect.Type) TypeSchemaValidator {
	c.schema = schema
	c.goType = goType
	return c
}
func (c typeSchemaValidatorContext) Errors() []string {
	return *c.errors
}

func (c typeSchemaValidatorContext) Validate() error {
	if isAny(c.goType) {
		return nil
	}
	if utils.IsMultiType(c.goType) {
		if _, err := utils.ExtractMultiTypeTypes(c.goType); err != nil {
			c.err(err.Error())
		}
	}
	if c.goType.Kind() == reflect.Pointer && !utils.IsMultiType(c.goType) {
		return c.WithType(c.goType.Elem()).Validate()
	}

	// Test global schema validation properties
	c.validateSchemaAllOf()
	c.validateSchemaNot()
	c.matchAllSchemaValidator("oneOf", c.schema.OneOf)
	c.matchAllSchemaValidator("anyOf", c.schema.AnyOf)

	// Test specific schema types validations
	c.validateObjectSchema()
	c.validateArraySchema()
	c.validateStringSchema()
	c.validateBooleanSchema()
	c.validateNumberSchema()
	c.validateIntegerSchema()

	if len(*c.errors) > 0 {
		return fmt.Errorf("%w %s", ErrSchemaIncompatibleWithType, c.goType)
	}

	return nil
}
