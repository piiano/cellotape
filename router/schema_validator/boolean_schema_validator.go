package schema_validator

import (
	"fmt"
	"reflect"
)

func (c typeSchemaValidatorContext) validateBooleanSchema() error {
	if c.schema.Type != booleanSchemaType {
		return nil
	}
	if c.goType.Kind() != reflect.Bool {
		return fmt.Errorf(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}
	return nil
}
