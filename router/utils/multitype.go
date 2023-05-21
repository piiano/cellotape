package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var ErrInvalidUseOfMultiType = errors.New("invalid use of MultiType")

var multiTypeReflectType = GetType[multiType]()

func IsMultiType(t reflect.Type) bool {
	return t.Implements(multiTypeReflectType)
}

func ExtractMultiTypeTypes(mtType reflect.Type) ([]reflect.Type, error) {
	multiTypeTypes := reflect.New(mtType).Elem().MethodByName("MultiTypeTypes")

	returnValues := multiTypeTypes.Call([]reflect.Value{})

	if err := returnValues[1].Interface(); err != nil {
		return nil, err.(error)
	}

	return returnValues[0].Interface().([]reflect.Type), nil
}

type multiType interface {
	MultiTypeTypes() ([]reflect.Type, error)
	json.Marshaler
	json.Unmarshaler
}

type MultiType[T any] struct {
	Values T
}

func (o *MultiType[T]) MultiTypeTypes() ([]reflect.Type, error) {
	fields, err := o.fields()
	if err != nil {
		return nil, err
	}

	return Map(fields, func(t reflect.StructField) reflect.Type {
		return t.Type
	}), nil
}

func (o *MultiType[T]) fields() ([]reflect.StructField, error) {
	structType := GetType[T]()
	if structType.Kind() != reflect.Struct {
		return []reflect.StructField{},
			fmt.Errorf("%w. expecting generic argument to be a struct",
				ErrInvalidUseOfMultiType)
	}
	fieldsMap := StructKeys(structType, "")
	for _, field := range fieldsMap {
		if field.Type.Kind() != reflect.Pointer {
			return []reflect.StructField{},
				fmt.Errorf("%w. field %q should be a pointer",
					ErrInvalidUseOfMultiType, field.Name)
		}
	}

	fields := Map(Entries(fieldsMap), func(e Entry[string, reflect.StructField]) reflect.StructField {
		return e.Value
	})

	if len(fields) == 0 {
		return nil, fmt.Errorf("%w. must have at least one field", ErrInvalidUseOfMultiType)
	}

	uniqueTypes := NewSet(Map(fields, func(t reflect.StructField) string {
		return t.Type.String()
	})...)

	if len(uniqueTypes) != len(fields) {
		return nil, fmt.Errorf("%w. each field of MultiType must be of a different type", ErrInvalidUseOfMultiType)
	}

	return fields, nil
}

func (o *MultiType[T]) MarshalJSON() ([]byte, error) {
	fields, err := o.fields()
	if err != nil {
		return nil, fmt.Errorf("can't marshal value to JSON due to %w", ErrInvalidUseOfMultiType)
	}

	value := reflect.ValueOf(o.Values)

	fields = Filter(fields, func(field reflect.StructField) bool {
		return !value.FieldByIndex(field.Index).IsNil()
	})

	if len(fields) == 0 {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(o),
			Str:   "non of MultiType fields is set",
		}
	}

	if len(fields) > 1 {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(o),
			Str:   "more than one field of MultiType is set",
		}
	}

	fieldValue := value.FieldByIndex(fields[0].Index)

	return json.Marshal(fieldValue.Interface())
}

func (o *MultiType[T]) UnmarshalJSON(bytes []byte) error {
	fields, err := o.fields()
	if err != nil {
		return fmt.Errorf("can't unmarshal JSON to value due to %w", ErrInvalidUseOfMultiType)
	}
	structValue := reflect.ValueOf(&o.Values).Elem()

	for _, field := range fields {
		fieldValue := structValue.FieldByIndex(field.Index)
		value := fieldValue.Addr().Interface()
		if err = json.Unmarshal(bytes, value); err == nil {
			return nil
		}
		fieldValue.Set(reflect.Zero(fieldValue.Type()))
	}

	unmarshalTypeError := &json.UnmarshalTypeError{}
	if ok := errors.As(err, &unmarshalTypeError); ok {
		unmarshalTypeError.Type = reflect.TypeOf(o)
		return unmarshalTypeError
	}

	return err
}
