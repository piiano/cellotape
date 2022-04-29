package schema_validator

import (
	"encoding/json"
	"fmt"
	"github.com/piiano/restcontroller/utils"
	"reflect"
	"strings"
)

func (c typeSchemaValidatorContext) validateObjectSchema() utils.MultiError {
	errs := utils.NewErrorsCollector()
	if SchemaType(c.schema.Type) != objectSchemaType {
		return nil
	}
	switch c.goType.Kind() {
	case reflect.Struct:
		if c.schema.Properties != nil {
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

// Extract the struct fields that are serializable as JSON
func structJsonFields(structType reflect.Type) map[string]reflect.StructField {
	visibleFields := utils.Filter(reflect.VisibleFields(structType), func(field reflect.StructField) bool {
		return !field.Anonymous
	})
	stringFields := make([]reflect.StructField, len(visibleFields))
	for i, field := range visibleFields {
		// change field type to string to prevent it from being omitted in json
		stringFields[i] = reflect.StructField{
			Type:      reflect.TypeOf(""),
			PkgPath:   field.PkgPath,
			Name:      field.Name,
			Tag:       field.Tag,
			Offset:    field.Offset,
			Index:     field.Index,
			Anonymous: field.Anonymous,
		}
	}
	//for i := 0; i < structType.NumField(); i++ {
	//	field := structType.Field(i)
	//}
	stringStructValue := reflect.New(reflect.StructOf(stringFields)).Elem()
	for _, field := range reflect.VisibleFields(stringStructValue.Type()) {
		fieldValue := stringStructValue.FieldByIndex(field.Index)
		if fieldValue.CanSet() {
			fieldValue.SetString("x")
		}
	}
	structJson, _ := json.Marshal(stringStructValue.Interface())
	mapValue := make(map[string]any)
	_ = json.Unmarshal(structJson, &mapValue)
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
