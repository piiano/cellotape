package schema_validator

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/piiano/cellotape/router/utils"
)

var textMarshallerType = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()

func (c typeSchemaValidatorContext) validateObjectSchema() {
	// TODO: validate required properties, nullable, additionalProperties, etc.
	serializedFromObject := isSerializedFromObject(c.goType)

	if !serializedFromObject {
		if c.schema.Type == openapi3.TypeObject {
			c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
		}
		return
	}

	if !isSchemaTypeObjectOrEmpty(c.schema) {
		c.err(schemaTypeIsIncompatibleWithType(c.schema, c.goType))
	}

	handleMultiType(func(t reflect.Type) bool {
		if t.Kind() == reflect.Struct {
			return c.assertStruct(t)
		}

		// kind must be struct or object because we validated above serializedFromObject
		return c.assertMap(t)
	})(c.goType)
}

func (c typeSchemaValidatorContext) assertStruct(t reflect.Type) bool {
	properties := c.schema.Properties
	if properties == nil {
		properties = make(map[string]*openapi3.SchemaRef, 0)
	}
	// TODO: add support for receiving in options how to extract type keys (to support schema for non-json serializers)
	fields := structJsonFields(t)

	validatedFields := utils.NewSet[string]()
	for name, field := range fields {

		validatedFields.Add(name)
		if property, ok := properties[name]; !ok {
			if additionalProperties := additionalPropertiesSchema(c.schema); additionalProperties == nil {
				c.err(fmt.Sprintf("field %q (%q) with type %s not found in object schema properties", field.Name, name, field.Type))
			} else if err := c.WithSchema(*additionalProperties).WithType(field.Type).Validate(); err != nil {
				c.err(fmt.Sprintf("field %q (%q) with type %s not found in object schema properties nor additonal properties", field.Name, name, field.Type))
			}
		} else if err := c.WithType(field.Type).WithSchema(*property.Value).Validate(); err != nil {
			c.err(schemaPropertyIsIncompatibleWithFieldType(name, field.Name, field.Type))
		}
	}
	for name := range properties {
		if !validatedFields.Has(name) {
			c.err(schemaPropertyIsNotMappedToFieldInType(name, t))
		}
	}

	return len(*c.errors) == 0
}
func (c typeSchemaValidatorContext) assertMap(t reflect.Type) bool {
	keyType := t.Key()
	mapValueType := t.Elem()

	// From the internal implementation of json.Marshal & json Unmarshal:
	// Map key must either have string kind, have an integer kind,
	// or be an encoding.TextUnmarshaler.
	switch keyType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		if !keyType.Implements(textMarshallerType) {
			c.err("object schema with map type must have a string compatible type. %s key is not string compatible", keyType)
		}
	}
	if c.schema.Properties != nil {

		for name, property := range c.schema.Properties {
			// check if property name is compatible with the map key type
			keyName, _ := json.Marshal(name)
			if err := json.Unmarshal(keyName, reflect.New(keyType).Interface()); err != nil {
				c.err("schema property name %q is incompatible with map key type %s", name, keyType)
			}
			// check if property schema is compatible with the map value type
			if err := c.WithType(mapValueType).WithSchema(*property.Value).Validate(); err != nil {
				c.err("schema property %q is incompatible with map value type %s", name, mapValueType)
			}
		}
	}

	return len(*c.errors) == 0
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

func additionalPropertiesSchema(schema openapi3.Schema) *openapi3.Schema {
	// if additional properties schema is defined explicitly return it
	if schema.AdditionalProperties != nil {
		return schema.AdditionalProperties.Value
	}

	// if additional properties is empty (tru by default) or set explicitly to true return an empty schema (schema for any type)
	if schema.AdditionalPropertiesAllowed == nil || *schema.AdditionalPropertiesAllowed {
		return openapi3.NewSchema()
	}

	// return nil if additional properties are not allowed
	return nil
}
