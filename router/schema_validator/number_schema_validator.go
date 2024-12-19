package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func (c typeSchemaValidatorContext) validateNumberSchema() {
	isGoTypeNumeric := isNumericType(c.goType)

	if !isGoTypeNumeric {
		if !c.schema.Type.Is(openapi3.TypeNumber) {
			return
		}

		// schema type is numeric and go type is not numeric
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
		return
	}

	// schema type is numeric and go type is not numeric
	if (c.schema.Format == floatFormat && !isFloat32(c.goType)) ||
		(c.schema.Format == doubleFormat && !isFloat64(c.goType)) {
		c.err(schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		return
	}

	// TODO: check type compatability with Max, ExclusiveMax, Min, and ExclusiveMin
}
