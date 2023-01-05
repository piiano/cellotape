package schema_validator

import "github.com/getkin/kin-openapi/openapi3"

func (c typeSchemaValidatorContext) validateBooleanSchema() {

	if isBoolType(c.goType) && !isSchemaTypeBooleanOrEmpty(c.schema) {
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}

	if c.schema.Type == openapi3.TypeBoolean && !isBoolType(c.goType) {
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}
}
