package schema_validator

import "github.com/getkin/kin-openapi/openapi3"

func (c typeSchemaValidatorContext) validateStringSchema() {
	if c.schema.Type.Is(openapi3.TypeString) && !isSerializedFromString(c.goType) {
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}

	// if schema type is not string and not empty
	if !isSchemaTypeStringOrEmpty(c.schema) {
		// and go type is a string
		if isString(c.goType) {
			// can't have a go type string for the remaining schema types: boolean, number, integer, array, object
			c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
		}

		// if schema type is not string and not empty other string validations has no meaning. return early.
		return
	}

	// if schema format is empty return early
	// pattern, minLength and maxLength have no effect on type validation.
	if c.schema.Format == "" {
		return
	}

	// if schema format is "byte" expect type to be compatible with []byte
	if c.schema.Format == byteFormat {
		if (c.schema.Type.Is(openapi3.TypeString) || isSerializedFromString(c.goType)) && !isSliceOfBytes(c.goType) {
			c.err(schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
		return
	}

	// if schema format is "uuid" expect type to be compatible with UUID
	if c.schema.Format == uuidFormat {
		if (c.schema.Type.Is(openapi3.TypeString) || isSerializedFromString(c.goType)) && !isUUIDCompatible(c.goType) {
			c.err(schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
		return
	}

	// if schema format is "date-time" or "time" expect type to be compatible with Time
	if isTimeFormat(c.schema) {
		if (c.schema.Type.Is(openapi3.TypeString) || isSerializedFromString(c.goType)) && !isTimeCompatible(c.goType) {
			c.err(schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		}
		return
	}

	// if schema format is any other string format expect go type to be string
	if isSchemaStringFormat(c.schema) && (c.schema.Type.Is(openapi3.TypeString) || isSerializedFromString(c.goType)) && !isString(c.goType) {
		c.err(schemaTypeWithFormatIsIncompatibleWithType(c.schema, c.goType))
		return
	}
}
