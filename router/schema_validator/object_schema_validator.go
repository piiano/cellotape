package schema_validator

import (
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/piiano/cellotape/router/utils"
	"reflect"
	"strings"
)

var textMarshallerType = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()

func (c typeSchemaValidatorContext) validateObjectSchema() error {
	// TODO: validate required properties, nullable, additionalProperties, etc.
	l := c.newLogger()
	if c.schema.Type != objectSchemaType {
		return nil
	}
	switch c.goType.Kind() {
	case reflect.Struct:
		if c.schema.Properties != nil {
			// TODO: add support for receiving in options how to extract type keys (to support schema for non-json serializers)
			fields := structJsonFields(c.goType)
			//c.schema.AdditionalPropertiesAllowed
			for name, field := range fields {
				property, ok := c.schema.Properties[name]
				if !ok {
					if c.schema.AdditionalPropertiesAllowed == nil || !*c.schema.AdditionalPropertiesAllowed {
						l.Logf(c.level, fmt.Sprintf("field %q (%q) with type %s not found in object schema properties", field.Name, name, field.Type))
						continue
					}
					if c.schema.AdditionalProperties.Value != nil {
						if err := c.WithSchema(*c.schema.AdditionalProperties.Value).WithType(field.Type).Validate(); err != nil {
							l.Logf(c.level, fmt.Sprintf("field %q (%q) with type %s not found in object schema properties nor additonal properties", field.Name, name, field.Type))
						}
					}
					continue
				}
				if err := c.WithType(field.Type).WithSchema(*property.Value).Validate(); err != nil {
					l.Logf(c.level, schemaPropertyIsIncompatibleWithFieldType(name, field.Name, field.Type))
				}
			}
			for name, property := range c.schema.Properties {
				field, ok := fields[name]
				if !ok {
					l.Logf(c.level, schemaPropertyIsNotMappedToFieldInType(name, c.goType))
					continue
				}
				if err := c.WithType(field.Type).WithSchema(*property.Value).Validate(); err != nil {
					l.Logf(c.level, schemaPropertyIsIncompatibleWithFieldType(name, field.Name, field.Type))
				}
			}
		}
	case reflect.Map:
		keyType := c.goType.Key()
		mapValueType := c.goType.Elem()
		switch keyType.Kind() {
		case reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		default:
			if !keyType.Implements(textMarshallerType) {
				l.Logf(c.level, "object schema with map type must have a string compatible type. %s key is not string compatible", keyType)
			}
		}
		if c.schema.Properties != nil {
			for name, property := range c.schema.Properties {
				// check if property name is compatible with the map key type
				keyName, _ := json.Marshal(name)
				if err := json.Unmarshal(keyName, reflect.New(keyType).Interface()); err != nil {
					l.Logf(c.level, "schema property name %q is incompatible with map key type %s", name, keyType)
				}
				// check if property schema is compatible with the map value type
				if err := c.WithType(mapValueType).WithSchema(*property.Value).Validate(); err != nil {
					l.Logf(c.level, "schema property %q is incompatible with map value type %s", name, mapValueType)
				}
			}
		}
	default:
		l.Logf(c.level, "object schema must be a struct type or a map. %s type is incompatible", c.goType)
	}
	return formatMustHaveNoError(l.MustHaveNoErrors(), c.schema.Type, c.goType)
}

// structJsonFields Extract the struct fields that are serializable as JSON corresponding to their JSON key
func structJsonFields(structType reflect.Type) map[string]reflect.StructField {
	// this method cares only about the json keys of the current struct in following json/encoding rules.
	// for simplification, we marshal an instance of the struct to json and unmarshal it to back to a map to get all json keys.
	//
	// To avoid fields filtering by omitempty tag we use bool type and set them with reflection to a non nil value.
	// this guarantee that we don't miss fields annotated with omitempty tag.
	// the original field types can be later extracted from the original structType

	// get all visible fields, embedded structs are excluded but their fields are included
	visibleFields := utils.Filter(reflect.VisibleFields(structType), func(field reflect.StructField) bool {
		return !field.Anonymous
	})
	// change field type to bool type
	boolFields := utils.Map(visibleFields, func(field reflect.StructField) reflect.StructField {
		field.Type = reflect.TypeOf(true)
		return field
	})
	//instantiate the struct type and set field values to true to prevent them from being omitted by omitempty tag
	boolStructValue := reflect.New(reflect.StructOf(boolFields)).Elem()
	for _, field := range reflect.VisibleFields(boolStructValue.Type()) {
		fieldValue := boolStructValue.FieldByIndex(field.Index)
		if fieldValue.CanSet() {
			fieldValue.SetBool(true)
		}
	}
	// marshal the struct to JSON
	structJson, _ := json.Marshal(boolStructValue.Interface())
	mapValue := make(map[string]bool)
	// unmarshal the JSON to map of bool
	_ = json.Unmarshal(structJson, &mapValue)
	// find the original StructField for each JSON key
	fields := make(map[string]reflect.StructField, len(mapValue))
	for key := range mapValue {
		fields[key], _ = structType.FieldByNameFunc(func(name string) bool {
			field, _ := structType.FieldByName(name)
			tag, found := field.Tag.Lookup("json")
			if !found {
				return name == key
			}
			name, _, _ = strings.Cut(tag, ",")
			return name == key
		})
	}
	return fields
}
