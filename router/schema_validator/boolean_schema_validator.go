package schema_validator

import "github.com/getkin/kin-openapi/openapi3"

func (c typeSchemaValidatorContext) validateBooleanSchema() {

	isTypeBool := isBoolType(c.goType)
	if (isTypeBool && !isSchemaTypeBooleanOrEmpty(c.schema)) ||
		(c.schema.Type.Is(openapi3.TypeBoolean) && !isTypeBool) {
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}
}
