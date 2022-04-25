package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router"
	"reflect"
	"testing"
)

func TestBooleanSchemaValidatorPassForBoolType(t *testing.T) {
	booleanSchema := openapi3.NewBoolSchema()
	booleanSchema.Title = "BoolSchema"
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *booleanSchema, router.Options{})
	boolType := reflect.TypeOf(true)
	if err := validator.WithType(boolType).Validate(); err != nil {
		t.Errorf("expect boolean schema to be compatible with %s type", boolType)
		t.Error(err)
	}
}

func TestBooleanSchemaValidatorFailOnWrongType(t *testing.T) {
	booleanSchema := openapi3.NewBoolSchema()
	booleanSchema.Title = "BoolSchema"
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *booleanSchema, router.Options{})
	boolArrayType := reflect.TypeOf([1]bool{true})
	if err := validator.WithType(boolArrayType).Validate(); err == nil {
		t.Errorf("expect boolean schema to be incompatible with %s type", boolArrayType)
	}
	intType := reflect.TypeOf(1)
	if err := validator.WithType(intType).Validate(); err == nil {
		t.Errorf("expect boolean schema to be incompatible with %s type", intType)
	}
}
