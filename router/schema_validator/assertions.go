package schema_validator

import (
	"reflect"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"

	"github.com/piiano/cellotape/router/utils"
)

var (
	timeType         = utils.GetType[time.Time]()
	uuidType         = utils.GetType[uuid.UUID]()
	sliceOfBytesType = utils.GetType[[]byte]()

	isString               = kindIs(reflect.String)
	isUUIDCompatible       = anyOf(isString, convertibleTo(uuidType))
	isSliceOfBytes         = anyOf(isString, typeIs(sliceOfBytesType))
	isTimeCompatible       = anyOf(isString, convertibleTo(timeType))
	isSerializedFromString = anyOf(isString, isUUIDCompatible, isTimeCompatible, isSliceOfBytes)

	isTimeFormat         = schemaFormatIs(dateTimeFormat, timeFormat)
	isSchemaStringFormat = schemaFormatIs(uuidFormat, byteFormat, dateTimeFormat, timeFormat, dateFormat, durationFormat,
		emailFormat, idnEmailFormat, hostnameFormat, idnHostnameFormat, ipv4Format, ipv6Format, uriFormat,
		uriReferenceFormat, iriFormat, iriReferenceFormat, uriTemplateFormat, jsonPointerFormat,
		relativeJsonPointerFormat, regexFormat, passwordFormat)

	isSerializedFromObject = allOf(kindIs(reflect.Struct, reflect.Map), not(isTimeCompatible))

	isSchemaTypeStringOrEmpty  = schemaTypeIsOrNil(openapi3.TypeString)
	isSchemaTypeBooleanOrEmpty = schemaTypeIsOrNil(openapi3.TypeBoolean)
	isSchemaTypeObjectOrEmpty  = schemaTypeIsOrNil(openapi3.TypeObject)
	isSchemaTypeArrayOrEmpty   = schemaTypeIsOrNil(openapi3.TypeArray)

	isBoolType    = kindIs(reflect.Bool)
	isFloat32     = kindIs(reflect.Float32)
	isFloat64     = kindIs(reflect.Float64)
	isNumericType = kindIs(reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64)

	isArrayGoType = allOf(kindIs(reflect.Array, reflect.Slice), not(isUUIDCompatible))
)

func isAny(t reflect.Type) bool {
	return t.Kind() == reflect.Interface && t.NumMethod() == 0
}

type assertion[T any] func(t T) bool
type typeAssertion = assertion[reflect.Type]
type schemaAssertion = assertion[openapi3.Schema]

func anyOf[A assertion[T], T any](assertions ...A) A {
	return func(t T) bool {
		for _, assert := range assertions {
			if assert(t) {
				return true
			}
		}
		return false
	}
}

func allOf[T any, A assertion[T]](assertions ...A) A {
	return func(t T) bool {
		for _, assert := range assertions {
			if !assert(t) {
				return false
			}
		}
		return true
	}
}

func not[T any, A assertion[T]](assertion A) A {
	return func(t T) bool {
		return !assertion(t)
	}
}

func schemaTypeIsOrNil(types ...string) schemaAssertion {
	return func(s openapi3.Schema) bool {
		for _, t := range types {
			if s.Type == nil || s.Type.Is(t) {
				return true
			}
		}
		return false
	}
}

func schemaFormatIs(types ...string) schemaAssertion {
	set := utils.NewSet(types...)
	return func(s openapi3.Schema) bool {
		return set.Has(s.Format)
	}
}

func kindIs(kinds ...reflect.Kind) typeAssertion {
	set := utils.NewSet(kinds...)
	return handleMultiType(func(t reflect.Type) bool {
		return set.Has(t.Kind())
	})
}

func typeIs(types ...reflect.Type) typeAssertion {
	set := utils.NewSet(types...)
	return handleMultiType(func(t reflect.Type) bool {
		return set.Has(t)
	})
}

func convertibleTo(targets ...reflect.Type) typeAssertion {
	return handleMultiType(func(t reflect.Type) bool {
		for _, target := range targets {
			if target.ConvertibleTo(t) {
				return true
			}
		}
		return false
	})
}

func handleMultiType(assertion typeAssertion) typeAssertion {
	return func(t reflect.Type) bool {
		if !utils.IsMultiType(t) {
			if t.Kind() == reflect.Pointer {
				return assertion(t.Elem())
			} else {
				return assertion(t)
			}
		}

		types, err := utils.ExtractMultiTypeTypes(t)
		if err != nil {
			return false
		}

		for _, mtType := range types {
			if assertion(mtType.Elem()) {
				return true
			}
		}

		return false
	}
}
