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

	isSchemaTypeStringOrEmpty  = schemaTypeIs(openapi3.TypeString, "")
	isSchemaTypeBooleanOrEmpty = schemaTypeIs(openapi3.TypeBoolean, "")
	isSchemaTypeObjectOrEmpty  = schemaTypeIs(openapi3.TypeObject, "")
	isSchemaTypeArrayOrEmpty   = schemaTypeIs(openapi3.TypeArray, "")
	isSchemaTypeNumberOrEmpty  = schemaTypeIs(openapi3.TypeNumber, "")

	isBoolType    = kindIs(reflect.Bool)
	isInt32       = kindIs(reflect.Int32)
	isInt64       = kindIs(reflect.Int64)
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

func schemaTypeIs(types ...string) schemaAssertion {
	return func(s openapi3.Schema) bool {
		return utils.NewSet(types...).Has(s.Type)
	}
}

func schemaFormatIs(types ...string) schemaAssertion {
	return func(s openapi3.Schema) bool {
		return utils.NewSet(types...).Has(s.Format)
	}
}

func kindIs(kinds ...reflect.Kind) typeAssertion {
	return handleMultiType(func(t reflect.Type) bool {
		return utils.NewSet(kinds...).Has(t.Kind())
	})
}

func typeIs(types ...reflect.Type) typeAssertion {
	return handleMultiType(func(t reflect.Type) bool {
		return utils.NewSet(types...).Has(t)
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
