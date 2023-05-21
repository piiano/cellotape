package schema_validator

import (
	"reflect"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestStringSchemaValidatorWithByteFormat(t *testing.T) {
	stringSchema := openapi3.NewStringSchema()
	stringSchema.Format = "byte"

	validator := schemaValidator(*stringSchema)
	errTemplate := "expect string schema to be compatible with %s type"
	bytes := reflect.TypeOf([]byte{})
	expectTypeToBeCompatible(t, validator, bytes, errTemplate, bytes)
}

func TestStringSchemaValidatorPassForStringType(t *testing.T) {
	stringSchema := openapi3.NewStringSchema()
	validator := schemaValidator(*stringSchema)
	errTemplate := "expect string schema to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, stringType, errTemplate, stringType)
}

// according to the spec the string validation properties should apply only when the type is set to string
func TestStringSchemaValidatorWithUntypedSchema(t *testing.T) {
	untypedSchemaWithUUIDFormat := openapi3.NewSchema().WithFormat(uuidFormat)

	otherNonStringTypes := utils.Filter(types, func(t reflect.Type) bool {
		return t != sliceOfBytesType && t != timeType
	})

	for _, validType := range otherNonStringTypes {
		t.Run(validType.String(), func(t *testing.T) {
			err := schemaValidator(*untypedSchemaWithUUIDFormat).WithType(validType).Validate()
			require.NoErrorf(t, err, "expect untyped schema to be compatible with %s type", validType)
		})
	}
}

func TestStringSchemaValidatorFailOnWrongType(t *testing.T) {
	stringSchema := openapi3.NewStringSchema()
	validator := schemaValidator(*stringSchema)
	errTemplate := "expect string schema to be incompatible with %s type"
	// filter string type from all defined test types
	var nonStringTypes = utils.Filter(types, func(t reflect.Type) bool {
		return t != stringType && t != reflect.PointerTo(stringType)
	})
	for _, nonStringType := range nonStringTypes {
		t.Run(nonStringType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonStringType, errTemplate, nonStringType)
		})
	}
}

func TestUUIDFormatSchemaValidator(t *testing.T) {
	uuidSchema := openapi3.NewStringSchema().WithFormat(uuidFormat)
	validator := schemaValidator(*uuidSchema)
	errTemplate := "expect string schema with uuid format to be %s with %s type"
	uuidType := reflect.TypeOf(uuid.New())
	expectTypeToBeCompatible(t, validator, uuidType, errTemplate, "compatible", uuidType)
	expectTypeToBeCompatible(t, validator, stringType, errTemplate, "compatible", stringType)
	// omit the uuid compatible types from all defined test types
	var nonUUIDCompatibleTypes = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		return t != uuidType && t != stringType && t != reflect.PointerTo(uuidType) && t != reflect.PointerTo(stringType)
	})
	for _, nonUUIDCompatibleType := range nonUUIDCompatibleTypes {
		t.Run(nonUUIDCompatibleType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonUUIDCompatibleType, errTemplate, "incompatible", nonUUIDCompatibleType)
		})
	}
}

func TestTimeFormatSchemaValidator(t *testing.T) {
	timeSchema := openapi3.NewStringSchema().WithFormat(timeFormat)

	errTemplate := "expect string schema with time format to be %s with %s type"
	timeCompatibleType := utils.NewSet(utils.Map([]reflect.Type{
		timeType, stringType, reflect.PointerTo(timeType), reflect.PointerTo(stringType),
	}, reflect.Type.String)...)

	for _, goType := range append(types, timeType, reflect.PointerTo(timeType)) {
		t.Run(goType.String(), func(t *testing.T) {
			err := schemaValidator(*timeSchema).WithType(goType).Validate()
			valid := timeCompatibleType.Has(goType.String())
			if valid {
				require.NoErrorf(t, err, errTemplate, "compatible", goType)
			} else {
				require.Errorf(t, err, errTemplate, "incompatible", goType)
			}
		})
	}
}

func TestDateTimeFormatSchemaValidator(t *testing.T) {
	timeSchema := openapi3.NewStringSchema().WithFormat(dateTimeFormat)

	errTemplate := "expect string schema with time format to be %s with %s type"
	timeType := reflect.TypeOf(time.Now())
	expectTypeToBeCompatible(t, schemaValidator(*timeSchema), timeType, errTemplate, "compatible", timeType)
	expectTypeToBeCompatible(t, schemaValidator(*timeSchema), stringType, errTemplate, "compatible", stringType)
	// omit the uuid compatible types from all defined test types
	var nonTimeCompatibleTypes = utils.Filter(types, func(t reflect.Type) bool {
		return t != timeType && t != stringType && t != reflect.PointerTo(timeType) && t != reflect.PointerTo(stringType)
	})
	for _, nonTimeCompatibleType := range nonTimeCompatibleTypes {
		t.Run(nonTimeCompatibleType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, schemaValidator(*timeSchema), nonTimeCompatibleType, errTemplate, "incompatible", nonTimeCompatibleType)
		})
	}
}

func TestStringSchemaValidatorWithOtherFormats(t *testing.T) {
	stringSchema := openapi3.NewStringSchema().WithFormat(hostnameFormat)
	validator := schemaValidator(*stringSchema)
	errTemplate := "expect string schema with time format to be %s with %s type"
	expectTypeToBeCompatible(t, validator, stringType, errTemplate, "compatible", stringType)
	// omit the string type from all defined test types
	var nonStringTypes = utils.Filter(types, func(t reflect.Type) bool {
		return t != stringType && t != reflect.PointerTo(stringType)
	})
	for _, nonStringType := range nonStringTypes {
		t.Run(nonStringType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonStringType, errTemplate, "incompatible", nonStringType)
		})
	}
}
