package schema_validator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router"
	"github.com/piiano/restcontroller/utils"
	"reflect"
	"testing"
)

var identifiableSchema = (&openapi3.Schema{Type: string(objectSchemaType)}).
	WithProperty("id", openapi3.NewStringSchema())
var personSchema = (&openapi3.Schema{Type: string(objectSchemaType)}).
	WithProperty("name", openapi3.NewStringSchema())
var identifiablePersonSchema = &openapi3.Schema{
	AllOf: openapi3.SchemaRefs{
		identifiableSchema.NewRef(),
		personSchema.NewRef(),
	},
}

type Identifiable struct {
	ID string `json:"id"`
}
type Person struct {
	Name string `json:"name"`
}
type IdentifiablePerson struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var identifiableType = reflect.TypeOf(Identifiable{ID: ""})
var personType = reflect.TypeOf(Person{Name: ""})
var identifiablePersonType = reflect.TypeOf(IdentifiablePerson{ID: "", Name: ""})
var stringToStringMapType = reflect.TypeOf(map[string]string{})

func TestSchemaAllOfValidatorPass(t *testing.T) {
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *identifiablePersonSchema, router.Options{})
	errTemplate := "expect schema with allOf property to be compatible with %s type"
	expectTypeToBeCompatible(t, validator, identifiablePersonType, errTemplate, identifiablePersonType)
	expectTypeToBeCompatible(t, validator, stringToStringMapType, errTemplate, stringToStringMapType)
}

func TestSchemaAllOfValidatorFailOnPartialMatchedType(t *testing.T) {
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *identifiablePersonSchema, router.Options{})
	errTemplate := "expect schema with allOf property to be incompatible with %s type"
	expectTypeToBeIncompatible(t, validator, identifiableType, errTemplate, identifiableType)
	expectTypeToBeIncompatible(t, validator, personType, errTemplate, personType)
}

func TestSchemaAllOfValidatorFailOnNoMatchedType(t *testing.T) {
	validator := NewTypeSchemaValidator(reflect.TypeOf(nil), *identifiablePersonSchema, router.Options{})
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
