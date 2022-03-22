package schema_validator

import (
	"fmt"
	"reflect"
)

func (c typeSchemaValidatorContext) validateArraySchema() error {
	if c.schema.Type != arraySchemaType {
		return nil
	}
	if c.goType.Kind() != reflect.Array && c.goType.Kind() != reflect.Slice {
		return fmt.Errorf(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}
	return c.WithSchemaAndType(*c.schema.Items.Value, c.goType.Elem()).Validate()
}
