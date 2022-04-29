package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/utils"
	"reflect"
	"testing"
)

type (
	Identifiable struct {
		ID string `json:"id"`
	}
	Person struct {
		Name string `json:"name"`
	}
	IdentifiablePerson struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

var (
	identifiableSchema       = openapi3.NewObjectSchema().WithProperty("id", openapi3.NewStringSchema())
	personSchema             = openapi3.NewObjectSchema().WithProperty("name", openapi3.NewStringSchema())
	identifiablePersonSchema = &openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			identifiableSchema.NewRef(),
			personSchema.NewRef(),
		},
	}
	identifiableType       = reflect.TypeOf(Identifiable{ID: ""})
	personType             = reflect.TypeOf(Person{Name: ""})
	identifiablePersonType = reflect.TypeOf(IdentifiablePerson{ID: "", Name: ""})
	stringToStringMapType  = reflect.TypeOf(map[string]string{})
)

func TestSchemaAllOfValidatorPass(t *testing.T) {
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *identifiablePersonSchema, Options{})
	errTemplate := "expect schema with allOf property to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, identifiablePersonType, errTemplate, identifiablePersonType)
	expectTypeToBeCompatible(t, validator, stringToStringMapType, errTemplate, stringToStringMapType)
}

func TestSchemaAllOfValidatorFailOnPartialMatchedType(t *testing.T) {
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *identifiablePersonSchema, Options{})
	errTemplate := "expect schema with allOf property to be incompatible with %s type"
	expectTypeToBeIncompatible(t, validator, identifiableType, errTemplate, identifiableType)
	expectTypeToBeIncompatible(t, validator, personType, errTemplate, personType)
}

func TestSchemaAllOfValidatorFailOnNoMatchedType(t *testing.T) {
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *identifiablePersonSchema, Options{})
	errTemplate := "expect schema with allOf property to be incompatible with %s type"
	invalidTypes := utils.Filter(types, func(t reflect.Type) bool {
		return t != stringToStringMapType
	})
	for _, invalidType := range invalidTypes {
		t.Run(invalidType.String(), func(t *testing.T) {
			expectTypeToBeIncompatible(t, validator, invalidType, errTemplate, invalidType)
		})
	}
}
