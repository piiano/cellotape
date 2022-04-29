package schema_validator

import (
	"encoding/json"
	"fmt"
	"github.com/piiano/restcontroller/utils"
	"reflect"
	"strings"
)

func (c typeSchemaValidatorContext) validateObjectSchema() utils.MultiError {
	// TODO: validate required properties, nullable, additionalProperties, etc.
	errs := utils.NewErrorsCollector()
	if SchemaType(c.schema.Type) != objectSchemaType {
		return nil
	}
	switch c.goType.Kind() {
	case reflect.Struct:
		if c.schema.Properties != nil {
			// TODO: add support for receiving in options how to extract type keys (to support schema for non-json serializers)
			fields := structJsonFields(c.goType)
			for name, property := range c.schema.Properties {
				field, ok := fields[name]
				if !ok {
					errs.AddIfNotNil(fmt.Errorf("property %q is not maped to a field in type %s", name, c.goType))
				}
				if ok {
					errs.AddIfNotNil(c.WithType(field.Type).WithSchema(*property.Value).Validate())
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
				errs.AddIfNotNil(fmt.Errorf("object schema with map type must have a string compatible type. %s key is not string compatible", keyType))
			}
		}
		if c.schema.Properties != nil {
			for name, property := range c.schema.Properties {
				// check if property name is compatible with the map key type
				keyName, _ := json.Marshal(name)
				if err := json.Unmarshal(keyName, reflect.New(keyType).Interface()); err != nil {
					errs.AddIfNotNil(fmt.Errorf("schema property name %q is incompatible with map key type %s", name, keyType))
				}
				// check if property schema is compatible with the map value type
				if err := c.WithType(mapValueType).WithSchema(*property.Value).Validate(); err != nil {
					errs.AddIfNotNil(fmt.Errorf("schema property %q is incompatible with map value type %s", name, mapValueType))
					errs.AddIfNotNil(err)
				}
			}
		}
	default:
		errs.AddIfNotNil(fmt.Errorf("object schema must be a struct type or a map. %s type is not compatible", c.goType))
	}
	return errs.ErrorOrNil()
}

// structJsonFields Extract the struct fields that are serializable as JSON corresponding to their JSON key
func structJsonFields(structType reflect.Type) map[string]reflect.StructField {
	// this method care only about the json keys of the current struct in following json/encoding rules.
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
