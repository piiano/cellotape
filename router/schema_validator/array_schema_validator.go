package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func (c typeSchemaValidatorContext) validateArraySchema() {
	if c.schema.Type == openapi3.TypeArray && !isArrayGoType(c.goType) {
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}

	if !isSchemaTypeArrayOrEmpty(c.schema) {
		if isArrayGoType(c.goType) && !isSliceOfBytes(c.goType) {
			c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
		}
		return
	}

	if isArrayGoType(c.goType) && c.schema.Items != nil {
		_ = c.WithSchemaAndType(*c.schema.Items.Value, c.goType.Elem()).Validate()
	}
}
