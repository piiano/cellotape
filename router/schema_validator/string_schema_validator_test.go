package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/piiano/cellotape/router/utils"
	"reflect"
	"testing"
	"time"
)

func TestStringSchemaValidatorPassForStringType(t *testing.T) {
	stringSchema := openapi3.NewStringSchema()
	validator := schemaValidator(*stringSchema)
	errTemplate := "expect string schema to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, stringType, errTemplate, stringType)
}

// according to the spec the string validation properties should apply only when the type is set to string
func TestStringSchemaValidatorWithUntypedSchema(t *testing.T) {
	untypedSchemaWithUUIDFormat := openapi3.NewSchema().WithFormat(uuidFormat)
	validator := schemaValidator(*untypedSchemaWithUUIDFormat)
	for _, validType := range types {
		t.Run(validType.String(), func(t *testing.T) {
			if err := validator.WithType(validType).validateStringSchema(); err != nil {
				t.Errorf("expect untyped schema to be compatible with %s type", validType)
			}
		})
	}
}

func TestStringSchemaValidatorFailOnWrongType(t *testing.T) {
	stringSchema := openapi3.NewStringSchema()
	validator := schemaValidator(*stringSchema)
	errTemplate := "expect string schema to be incompatible with %s type"
	// filter string type from all defined test types
	var nonStringTypes = utils.Filter(types, func(t reflect.Type) bool {
		return t != stringType
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
		return t != uuidType && t != stringType
	})
	for _, nonUUIDCompatibleType := range nonUUIDCompatibleTypes {
		t.Run(nonUUIDCompatibleType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonUUIDCompatibleType, errTemplate, "incompatible", nonUUIDCompatibleType)
		})
	}
}

func TestTimeFormatSchemaValidator(t *testing.T) {
	timeSchema := openapi3.NewStringSchema().WithFormat(timeFormat)
	validator := schemaValidator(*timeSchema)
	errTemplate := "expect string schema with time format to be %s with %s type"
	timeType := reflect.TypeOf(time.Now())
	expectTypeToBeCompatible(t, validator, timeType, errTemplate, "compatible", timeType)
	expectTypeToBeCompatible(t, validator, stringType, errTemplate, "compatible", stringType)
	// omit the uuid compatible types from all defined test types
	var nonTimeCompatibleTypes = utils.Filter(types, func(t reflect.Type) bool {
		return t != timeType && t != stringType
	})
	for _, nonTimeCompatibleType := range nonTimeCompatibleTypes {
		t.Run(nonTimeCompatibleType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonTimeCompatibleType, errTemplate, "incompatible", nonTimeCompatibleType)
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
		return t != stringType
	})
	for _, nonStringType := range nonStringTypes {
		t.Run(nonStringType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, nonStringType, errTemplate, "incompatible", nonStringType)
		})
	}
}
