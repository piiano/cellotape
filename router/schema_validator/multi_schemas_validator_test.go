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

type composer func(...*openapi3.Schema) *openapi3.Schema

type multiSchemaCase struct {
	name string
	composer
}

var multiSchemaCases = []multiSchemaCase{
	{name: "oneOf", composer: openapi3.NewOneOfSchema},
	{name: "anyOf", composer: openapi3.NewAnyOfSchema},
}

func TestMultiSchemaValidator(t *testing.T) {
	for _, schemaCase := range multiSchemaCases {
		t.Run(schemaCase.name, func(t *testing.T) {
			testCases := []struct {
				name         string
				goType       reflect.Type
				schema       *openapi3.Schema
				errAssertion func(require.TestingT, error, ...any)
			}{
				{
					name:   "multiple different types",
					goType: utils.GetType[any](),
					schema: schemaCase.composer(
						openapi3.NewBoolSchema(),
						openapi3.NewStringSchema(),
						openapi3.NewInt64Schema(),
					),
					errAssertion: require.NoError,
				},
				{
					name: "multiple different types",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						A *bool
						B *string
						C *int64
					}]{}),
					schema: schemaCase.composer(
						openapi3.NewBoolSchema(),
						openapi3.NewStringSchema(),
						openapi3.NewInt64Schema(),
					),
					errAssertion: require.NoError,
				},
				{
					name:   "simple go type for schema with multiple different types",
					goType: boolType,
					schema: schemaCase.composer(
						openapi3.NewBoolSchema(),
						openapi3.NewStringSchema(),
						openapi3.NewInt64Schema(),
					),
					errAssertion: require.Error,
				},
				{
					name: "missing go type for schema with multiple different types",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						A *bool
						B *string
					}]{}),
					schema: schemaCase.composer(
						openapi3.NewBoolSchema(),
						openapi3.NewStringSchema(),
						openapi3.NewInt64Schema(),
					),
					errAssertion: require.Error,
				},
				{
					name: "missing schema type for multi type go type",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						A *bool
						B *string
						C *int64
					}]{}),
					schema: schemaCase.composer(
						openapi3.NewBoolSchema(),
						openapi3.NewStringSchema(),
					),
					errAssertion: require.Error,
				},
				{
					name: "insignificant of types and schema order",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						B *string
						A *bool
					}]{}),
					schema: schemaCase.composer(
						openapi3.NewBoolSchema(),
						openapi3.NewStringSchema(),
					),
					errAssertion: require.NoError,
				},
				{
					name:   "simple single type with different validations",
					goType: intType,
					schema: func() *openapi3.Schema {
						schema := schemaCase.composer(
							&openapi3.Schema{MultipleOf: utils.Ptr(3.0)},
							&openapi3.Schema{MultipleOf: utils.Ptr(5.0)},
						)
						schema.Type = openapi3.TypeNumber
						return schema
					}(),
					errAssertion: require.NoError,
				},
				{
					name: "simple uuid or time",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						A *uuid.UUID
						B *time.Time
					}]{}),
					schema: func() *openapi3.Schema {
						schema := schemaCase.composer(
							&openapi3.Schema{Format: "uuid"},
							&openapi3.Schema{Format: "date-time"},
						)
						schema.Type = openapi3.TypeString
						return schema
					}(),
					errAssertion: require.NoError,
				},
				{
					name:   "integer formats",
					goType: int64Type,
					schema: func() *openapi3.Schema {
						schema := schemaCase.composer(
							&openapi3.Schema{Format: int64Format},
							&openapi3.Schema{Format: int32Format},
						)
						schema.Type = openapi3.TypeInteger
						return schema
					}(),
					errAssertion: require.NoError,
				},
				{
					name: "one struct or another",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						A *struct{ A string }
						B *struct{ B string }
					}]{}),
					schema: func() *openapi3.Schema {
						schema := schemaCase.composer(
							&openapi3.Schema{Properties: openapi3.Schemas{"A": openapi3.NewStringSchema().NewRef()}},
							&openapi3.Schema{Properties: openapi3.Schemas{"B": openapi3.NewStringSchema().NewRef()}},
						)
						schema.Type = openapi3.TypeObject
						return schema
					}(),
					errAssertion: require.NoError,
				},
				{
					name: "one struct or another with error",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						A *struct{ A string }
						B *bool
					}]{}),
					schema: func() *openapi3.Schema {
						schema := schemaCase.composer(
							&openapi3.Schema{Properties: openapi3.Schemas{"A": openapi3.NewStringSchema().NewRef()}},
							&openapi3.Schema{Properties: openapi3.Schemas{"B": openapi3.NewStringSchema().NewRef()}},
						)
						schema.Type = openapi3.TypeObject
						return schema
					}(),
					errAssertion: require.Error,
				},
				{
					name: "one array or another",
					goType: reflect.TypeOf(&utils.MultiType[struct {
						A *[]string
						B *[]bool
					}]{}),
					schema: func() *openapi3.Schema {
						schema := schemaCase.composer(
							&openapi3.Schema{Items: openapi3.NewStringSchema().NewRef()},
							&openapi3.Schema{Items: openapi3.NewBoolSchema().NewRef()},
						)
						schema.Type = openapi3.TypeArray
						return schema
					}(),
					errAssertion: require.NoError,
				},
			}

			for _, testCase := range testCases {
				t.Run(testCase.name, func(t *testing.T) {
					err := schemaValidator(*testCase.schema).WithType(testCase.goType).Validate()
					testCase.errAssertion(t, err)
				})
			}
		})
	}
}

func TestSchemaMultiSchemaValidatorFailOnNoMatchedType(t *testing.T) {
	for _, schemaCase := range multiSchemaCases {
		t.Run(schemaCase.name, func(t *testing.T) {
			schema := schemaCase.composer(
				openapi3.NewBoolSchema(),
				openapi3.NewStringSchema(),
				openapi3.NewInt64Schema(),
			)
			validator := schemaValidator(*schema)

			errTemplate := "expect schema with %s property to be incompatible with %s type"
			for _, invalidType := range types {
				t.Run(invalidType.String(), func(t *testing.T) {
					expectTypeToBeIncompatible(t, validator, invalidType, errTemplate, schemaCase.name, invalidType)
				})
			}
		})
	}
}

func TestCorruptedMultiType(t *testing.T) {
	testType := utils.GetType[*utils.MultiType[bool]]()

	isMultiType := utils.IsMultiType(testType)
	require.True(t, isMultiType)

	_, err := utils.ExtractMultiTypeTypes(testType)
	require.Error(t, err)

	err = schemaValidator(*openapi3.NewOneOfSchema(
		openapi3.NewStringSchema(),
		openapi3.NewBoolSchema(),
	)).WithType(testType).Validate()
	require.Error(t, err, "expect schema to be incompatible with invalid MultiType %s", testType)
}
