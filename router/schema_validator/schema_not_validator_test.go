package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/cellotape/router/utils"
	"reflect"
	"testing"
)

func TestSchemaNotValidatorPass(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	notBooleanSchema.Not = openapi3.NewBoolSchema().NewRef()
	validator := schemaValidator(*notBooleanSchema)
	// filter bool type from all defined test types
	var nonBoolTypes = utils.Filter[reflect.Type](types, func(t reflect.Type) bool {
		return t != boolType && t != reflect.PointerTo(boolType)
	})
	errTemplate := "expect schema with not bool schema to be compatible with %s type"
	for _, nonBoolType := range nonBoolTypes {
		t.Run(nonBoolType.String(), func(t *testing.T) {
			expectTypeToBeCompatible(t, validator, nonBoolType, errTemplate, nonBoolType)
		})
	}
}

func TestSchemaNotValidatorFailOnWrongType(t *testing.T) {
	notBooleanSchema := openapi3.NewSchema()
	notBooleanSchema.Not = openapi3.NewBoolSchema().NewRef()
	validator := schemaValidator(*notBooleanSchema)
	errTemplate := "expect schema with not bool schema to be incompatible with %s type"
	expectTypeToBeIncompatible(t, validator, boolType, errTemplate, boolType)
}
