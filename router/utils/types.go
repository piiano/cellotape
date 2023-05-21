package utils

import (
	"reflect"
	"strings"
)

// GetType returns reflect.Type of the generic parameter it receives.
func GetType[T any]() reflect.Type { return reflect.TypeOf(new(T)).Elem() }

// Nil represents an empty type.
// You can use it with the HandlerFunc generic parameters to declare no Request with no request body, no path or query
// params, or responses with no response body.
type Nil *uintptr

// NilType represent the type of Nil.
var NilType = GetType[Nil]()

const ignoreFieldTagValue = "-"

// StructKeys returns a map of "key" -> "field" for all the fields in the struct.
// Key is the field tag if exists or the field name otherwise.
// Field is the reflect.StructField of the field.
//
// Unexported fields or fields with tag value of "-" are ignored.
//
// StructKeys will recursively traverse all the embedded structs and return their fields as well.
func StructKeys(structType reflect.Type, tag string) map[string]reflect.StructField {
	if structType == nil || structType == NilType {
		return map[string]reflect.StructField{}
	}
	return FromEntries(ConcatSlices(Map(Filter(reflect.VisibleFields(structType), func(field reflect.StructField) bool {
		return !field.Anonymous && field.IsExported() && field.Tag.Get(tag) != ignoreFieldTagValue
	}), func(field reflect.StructField) Entry[string, reflect.StructField] {
		name := field.Tag.Get(tag)
		name, _, _ = strings.Cut(name, ",")
		if name == "" {
			name = field.Name
		}
		return Entry[string, reflect.StructField]{Key: name, Value: field}
	}), ConcatSlices(Map(Filter(reflect.VisibleFields(structType), func(field reflect.StructField) bool {
		return field.Anonymous && field.IsExported() && field.Tag.Get(tag) != ignoreFieldTagValue
	}), func(field reflect.StructField) []Entry[string, reflect.StructField] {
		return Entries(StructKeys(field.Type, tag))
	})...)))
}
