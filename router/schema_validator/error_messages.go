package schema_validator

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

func schemaAllOfPropertyIncompatibleWithType(invalidOptions int, options int, goType reflect.Type) string {
	return fmt.Sprintf("%d/%d schemas defined in allOf are incompatible with type %s", invalidOptions, options, goType)
}

func schemaTypeWithFormatIsIncompatibleWithType(schema openapi3.Schema, goType reflect.Type) string {
	return fmt.Sprintf("%s schema with %s format is incompatible with type %s", schema.Type, schema.Format, goType)
}

func schemaTypeIsIncompatibleWithType(schema openapi3.Schema, goType reflect.Type) string {
	return fmt.Sprintf("%s schema is incompatible with type %s", schema.Type, goType)
}

func schemaPropertyIsNotMappedToFieldInType(name string, fieldType reflect.Type) string {
	return fmt.Sprintf("property %q is not mapped to a field in type %s", name, fieldType)
}

func schemaPropertyIsIncompatibleWithFieldType(property string, field string, fieldType reflect.Type) string {
	return fmt.Sprintf("property %q is incompatible with type %s of field %q", property, fieldType, field)
}
