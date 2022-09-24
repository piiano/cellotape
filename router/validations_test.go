package router

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piiano/cellotape/router/utils"
)

func TestValidateContentTypes(t *testing.T) {
	err := validateContentTypes(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, utils.NewSet[string]())
	require.NoError(t, err)
}

func TestValidateContentTypesWithJSONContentType(t *testing.T) {
	err := validateContentTypes(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
		spec: OpenAPISpec(openapi3.T{
			Paths: openapi3.Paths{
				"/": &openapi3.PathItem{
					Get: &openapi3.Operation{
						RequestBody: &openapi3.RequestBodyRef{
							Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewSchema()),
						},
					},
				},
			},
		}),
	}, utils.NewSet[string]())
	require.NoError(t, err)
}

func TestValidateContentTypesWithExcludedOperation(t *testing.T) {
	err := validateContentTypes(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
		spec: OpenAPISpec(openapi3.T{
			Paths: openapi3.Paths{
				"/": &openapi3.PathItem{
					Get: &openapi3.Operation{
						OperationID: "foo",
						RequestBody: &openapi3.RequestBodyRef{
							Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewSchema()),
						},
					},
				},
			},
		}),
	}, utils.NewSet[string]("foo"))
	require.NoError(t, err)
}

func TestValidateContentTypesErrorWithMissingJSONContentType(t *testing.T) {
	err := validateContentTypes(openapi{
		options:      DefaultOptions(),
		contentTypes: ContentTypes{},
		spec: OpenAPISpec(openapi3.T{
			Paths: openapi3.Paths{
				"/": &openapi3.PathItem{
					Get: &openapi3.Operation{
						RequestBody: &openapi3.RequestBodyRef{
							Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewSchema()),
						},
					},
				},
			},
		}),
	}, utils.NewSet[string]())
	require.Error(t, err)
}

func TestValidateHandleAllPathParams(t *testing.T) {
	counter := validateHandleAllPathParams(openapi{
		options: DefaultOptions(),
	}, PropagateError, operation{
		handler: handler{
			request: requestTypes{
				pathParams: reflect.TypeOf(struct {
					Param string `uri:"foo"`
				}{}),
			},
		},
	}, SpecOperation{
		Operation: &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: openapi3.NewPathParameter("foo").WithSchema(openapi3.NewStringSchema()),
				},
				&openapi3.ParameterRef{
					Value: openapi3.NewPathParameter("bar").WithSchema(openapi3.NewStringSchema()),
				},
			},
		},
	})
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateHandleAllQueryParams(t *testing.T) {
	counter := validateHandleAllQueryParams(openapi{
		options: DefaultOptions(),
	}, PropagateError, operation{
		handler: handler{
			request: requestTypes{
				queryParams: reflect.TypeOf(struct {
					Param string `form:"foo"`
				}{}),
			},
		},
	}, SpecOperation{
		Operation: &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: openapi3.NewQueryParameter("foo").WithSchema(openapi3.NewStringSchema()),
				},
				&openapi3.ParameterRef{
					Value: openapi3.NewQueryParameter("bar").WithSchema(openapi3.NewStringSchema()),
				},
			},
		},
	})
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateHandleAllResponses(t *testing.T) {
	counter := validateHandleAllResponses(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, operation{
		handler: handler{
			responses: handlerResponses{
				200: httpResponse{
					status:       200,
					responseType: reflect.TypeOf(""),
					fieldIndex:   []int{0},
				},
			},
		},
	}, SpecOperation{
		Operation: &openapi3.Operation{
			Responses: testSpecResponse("200", "application/json", openapi3.NewStringSchema()),
		},
	})
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateHandleAllResponsesError(t *testing.T) {
	counter := validateHandleAllResponses(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, operation{
		handler: handler{
			responses: handlerResponses{
				200: httpResponse{
					status:       200,
					responseType: reflect.TypeOf(""),
					fieldIndex:   []int{0},
				},
			},
		},
	}, SpecOperation{
		Operation: &openapi3.Operation{
			Responses: testSpecResponse("200", "application/json", openapi3.NewStringSchema()),
		},
	})
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateHandleAllResponsesInvalidStatusError(t *testing.T) {
	counter := validateHandleAllResponses(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, operation{
		handler: handler{
			responses: handlerResponses{
				200: httpResponse{
					status:       200,
					responseType: reflect.TypeOf(""),
					fieldIndex:   []int{0},
				},
			},
		},
	}, SpecOperation{
		Operation: &openapi3.Operation{
			Responses: testSpecResponse("20x", "application/json", openapi3.NewStringSchema()),
		},
	})
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateHandleAllResponsesMissingStatusError(t *testing.T) {
	counter := validateHandleAllResponses(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, operation{
		handler: handler{
			responses: handlerResponses{},
		},
	}, SpecOperation{
		Operation: &openapi3.Operation{
			Responses: testSpecResponse("200", "application/json", openapi3.NewStringSchema()),
		},
	})
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateRequestBodyType(t *testing.T) {
	counter := validateRequestBodyType(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, handler{
		request: requestTypes{
			requestBody: reflect.TypeOf(""),
		},
	}, &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: openapi3.NewStringSchema(),
					},
				},
			},
		},
	}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateRequestBodyTypeIgnoreMissingContentType(t *testing.T) {
	counter := validateRequestBodyType(openapi{
		options: DefaultOptions(),
	}, PropagateError, handler{
		request: requestTypes{
			requestBody: reflect.TypeOf(""),
		},
	}, &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: openapi3.NewStringSchema(),
					},
				},
			},
		},
	}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateRequestBodyTypeErrorWithNoBodyInSpec(t *testing.T) {
	counter := validateRequestBodyType(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, handler{
		request: requestTypes{
			requestBody: reflect.TypeOf(""),
		},
	}, &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: openapi3.NewIntegerSchema(),
					},
				},
			},
		},
	}, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateRequestBodyTypeErrorWithIcompatibleSchema(t *testing.T) {
	counter := validateRequestBodyType(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, handler{
		request: requestTypes{
			requestBody: reflect.TypeOf(""),
		},
	}, nil, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateQueryParamsType(t *testing.T) {
	counter := validateQueryParamsType(openapi{}, PropagateError, handler{}, openapi3.Parameters{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
	counter = validateQueryParamsType(openapi{
		options: DefaultOptions(),
	}, PropagateError, handler{
		request: requestTypes{
			queryParams: reflect.TypeOf(struct {
				Param string `form:"foo"`
			}{}),
		},
	}, openapi3.Parameters{
		&openapi3.ParameterRef{
			Value: openapi3.NewQueryParameter("foo").WithSchema(openapi3.NewStringSchema()),
		},
	}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateCollidingEmbeddedQueryQueryParamsType(t *testing.T) {
	counter := validateQueryParamsType(openapi{}, PropagateError, handler{}, openapi3.Parameters{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
	counter = validateHandleAllQueryParams(openapi{
		options: DefaultOptions(),
	}, PropagateError, operation{
		handler: handler{
			request: requestTypes{
				queryParams: reflect.TypeOf(CollidingFieldsParams{}),
			},
		},
	}, SpecOperation{
		Operation: &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: openapi3.NewQueryParameter("param1").WithSchema(openapi3.NewStringSchema()),
				},
				&openapi3.ParameterRef{
					Value: openapi3.NewQueryParameter("param2").WithSchema(openapi3.NewStringSchema()),
				},
			},
		},
	})
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateQueryParamsTypeFailWhenMissingInSpec(t *testing.T) {
	counter := validateQueryParamsType(openapi{}, PropagateError, handler{}, openapi3.Parameters{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
	counter = validateQueryParamsType(openapi{
		options: DefaultOptions(),
	}, PropagateError, handler{
		request: requestTypes{
			queryParams: reflect.TypeOf(struct {
				Param string `form:"foo"`
			}{}),
		},
	}, openapi3.Parameters{}, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateQueryParamsTypeFailWhenIncompatibleType(t *testing.T) {
	counter := validateQueryParamsType(openapi{}, PropagateError, handler{}, openapi3.Parameters{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
	counter = validateQueryParamsType(openapi{
		options: DefaultOptions(),
	}, PropagateError, handler{
		request: requestTypes{
			queryParams: reflect.TypeOf(struct {
				Param string `form:"foo"`
			}{}),
		},
	}, openapi3.Parameters{
		&openapi3.ParameterRef{
			Value: openapi3.NewQueryParameter("foo").WithSchema(openapi3.NewIntegerSchema()),
		},
	}, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidatePathParamsType(t *testing.T) {
	counter := validatePathParamsType(openapi{}, PropagateError, handler{}, openapi3.Parameters{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
	counter = validatePathParamsType(openapi{
		options: DefaultOptions(),
	}, PropagateError, handler{
		request: requestTypes{
			pathParams: reflect.TypeOf(struct {
				Param string `uri:"foo"`
			}{}),
		},
	}, openapi3.Parameters{
		&openapi3.ParameterRef{
			Value: openapi3.NewPathParameter("foo").WithSchema(openapi3.NewStringSchema()),
		},
	}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidatePathParamsTypeFailWhenMissingInSpec(t *testing.T) {
	counter := validatePathParamsType(openapi{}, PropagateError, handler{}, openapi3.Parameters{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
	counter = validatePathParamsType(openapi{
		options: DefaultOptions(),
	}, PropagateError, handler{
		request: requestTypes{
			pathParams: reflect.TypeOf(struct {
				Param string `uri:"foo"`
			}{}),
		},
	}, openapi3.Parameters{}, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidatePathParamsTypeFailWhenIncompatibleType(t *testing.T) {
	counter := validatePathParamsType(openapi{}, PropagateError, handler{}, openapi3.Parameters{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
	counter = validatePathParamsType(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, handler{
		request: requestTypes{
			pathParams: reflect.TypeOf(struct {
				Param string `uri:"foo"`
			}{}),
		},
	}, openapi3.Parameters{
		&openapi3.ParameterRef{
			Value: openapi3.NewPathParameter("foo").WithSchema(openapi3.NewIntegerSchema()),
		},
	}, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestStructKeys(t *testing.T) {
	structType := reflect.TypeOf(struct {
		Field1 string `json:"field1"`
		Field2 int    `json:",omitempty"`
		Field3 bool
	}{})
	keys := structKeys(structType, "json")
	assert.Equal(t, map[string]reflect.StructField{
		"field1": structType.Field(0),
		"Field2": structType.Field(1),
		"Field3": structType.Field(2),
	}, keys)

	structType2 := reflect.TypeOf(struct {
		Field1 string `form:"field1"`
		Field2 int    `form:",omitempty"`
		Field3 bool
	}{})
	keys2 := structKeys(structType2, "form")
	assert.Equal(t, map[string]reflect.StructField{
		"field1": structType2.Field(0),
		"Field2": structType2.Field(1),
		"Field3": structType2.Field(2),
	}, keys2)
}

func TestValidateResponseTypes(t *testing.T) {
	counter := validateResponseTypes(openapi{}, PropagateError, handler{}, &openapi3.Operation{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)

	counter = validateResponseTypes(openapi{
		contentTypes: DefaultContentTypes(),
	}, PropagateError, handler{
		responses: handlerResponses{
			200: httpResponse{
				status:       200,
				responseType: reflect.TypeOf(""),
			},
		},
	}, &openapi3.Operation{
		Responses: testSpecResponse("200", "application/json", openapi3.NewStringSchema()),
	}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateResponseTypesIgnoreMissingContentType(t *testing.T) {
	counter := validateResponseTypes(openapi{}, PropagateError, handler{}, &openapi3.Operation{}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)

	counter = validateResponseTypes(openapi{}, PropagateError, handler{
		responses: handlerResponses{
			200: httpResponse{
				status:       200,
				responseType: reflect.TypeOf(""),
			},
		},
	}, &openapi3.Operation{
		Responses: testSpecResponse("200", "application/json", openapi3.NewStringSchema()),
	}, "")
	assert.Equal(t, 0, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateResponseTypesMissingStatusErr(t *testing.T) {
	counter := validateResponseTypes(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, handler{
		responses: handlerResponses{
			500: httpResponse{
				status:       500,
				responseType: reflect.TypeOf(""),
			},
		},
	}, &openapi3.Operation{
		Responses: testSpecResponse("200", "application/json", openapi3.NewStringSchema()),
	}, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func TestValidateResponseTypesIncompatibleTypeErr(t *testing.T) {
	counter := validateResponseTypes(openapi{
		options:      DefaultOptions(),
		contentTypes: DefaultContentTypes(),
	}, PropagateError, handler{
		responses: handlerResponses{
			200: httpResponse{
				status:       200,
				responseType: reflect.TypeOf(""),
			},
		},
	}, &openapi3.Operation{
		Responses: testSpecResponse("200", "application/json", openapi3.NewIntegerSchema()),
	}, "")
	assert.Equal(t, 1, counter.Errors)
	assert.Equal(t, 0, counter.Warnings)
}

func testSpecResponse(status string, contentType string, schema *openapi3.Schema) map[string]*openapi3.ResponseRef {
	return map[string]*openapi3.ResponseRef{
		status: {
			Value: &openapi3.Response{
				Content: openapi3.Content{
					contentType: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Value: schema,
						},
					},
				},
			},
		},
	}
}
