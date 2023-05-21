package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func (c typeSchemaValidatorContext) validateArraySchema() {
	isGoTypeArray := isArrayGoType(c.goType)
	if c.schema.Type == openapi3.TypeArray && !isGoTypeArray {
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}

	if !isSchemaTypeArrayOrEmpty(c.schema) {
		if isGoTypeArray && !isSliceOfBytes(c.goType) {
			c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
		}
		return
	}

	if isGoTypeArray && c.schema.Items != nil {
		_ = c.WithSchemaAndType(*c.schema.Items.Value, c.goType.Elem()).Validate()
	}
}
