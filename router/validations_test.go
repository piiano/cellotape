package router

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piiano/restcontroller/router/utils"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestValidateOpenAPIRouter(t *testing.T) {
	oa := testData()

	operations := flattenOperations(oa.group)
	err := validateOpenAPIRouter(&oa, operations)
	if err != nil {
		log.Println(err)
	}
}

func testData() openapi {
	oa := openapi{
		options: DefaultOptions(),
		spec: OpenAPISpec{
			Paths: openapi3.Paths{
				"/": {
					Post:   &openapi3.Operation{OperationID: "id1"},
					Get:    &openapi3.Operation{OperationID: "id2"},
					Put:    &openapi3.Operation{OperationID: "id3"},
					Delete: &openapi3.Operation{OperationID: "id4"},
				},
			},
		},
		group: group{
			groups: []group{
				{operations: []operation{
					{id: "id1", handler: handler{}},
					{id: "id2", handler: handler{}},
				}},
			},
			operations: []operation{
				{id: "id3", handler: handler{}},
				{id: "id4", handler: handler{}},
			},
		},
	}
	return oa
}

func testOperation() operation {

	return operation{id: "id3", handler: handler{}}
}

func TestValidateHandleAllOperations(t *testing.T) {}
func TestValidateContentTypes(t *testing.T)        {}
func TestValidateOperation(t *testing.T)           {}
func TestValidateHandleAllResponses(t *testing.T)  {}
func TestValidateRequestBodyType(t *testing.T)     {}
func TestValidatePathParamsType(t *testing.T)      {}
func TestValidateQueryParamsType(t *testing.T)     {}
func TestValidateParamsType(t *testing.T)          {}
func TestValidateResponseTypes(t *testing.T)       {}

func mockStructField(name string, t reflect.Type) reflect.StructField {
	return reflect.StructField{
		Name: strings.ToTitle(name),
		Type: nil,
		Tag:  reflect.StructTag(fmt.Sprintf("json:%q", strings.ToLower(name))),
	}
}

func mockStruct(fields map[string]reflect.Type) reflect.Type {
	return reflect.StructOf(utils.Map(utils.Entries(fields), func(entry utils.Entry[string, reflect.Type]) reflect.StructField {
		return mockStructField(entry.Key, entry.Value)
	}))
}
func mockParamsStruct(params []param) reflect.Type {
	return reflect.StructOf(utils.Map(params, func(p param) reflect.StructField {
		return mockStructField(p.name, p.paramType)
	}))
}

func mockPathParamsType() {
	mockStruct(map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
	})
}

type param struct {
	name      string
	paramType reflect.Type
}

var (
	intType      = reflect.TypeOf(0)
	stringType   = reflect.TypeOf("")
	typeToSchema = map[reflect.Type]*openapi3.Schema{
		intType:    openapi3.NewIntegerSchema(),
		stringType: openapi3.NewStringSchema(),
	}
)

func testCases() {
	params := map[string][]param{
		"no params":         {},
		"one string params": {{name: "param1", paramType: stringType}},
		"one int params":    {{name: "param1", paramType: intType}},
		"two params":        {{name: "param1", paramType: stringType}, {name: "param2", paramType: intType}},
	}
	utils.Map(utils.Entries(params), func(entry utils.Entry[string, []param]) reflect.Type {
		return mockParamsStruct(entry.Value)
	})
	utils.Map(utils.Entries(params), func(entry utils.Entry[string, []param]) openapi3.Parameters {
		return utils.Map(entry.Value, func(p param) *openapi3.ParameterRef {
			return &openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name: p.name,
					In:   "",
					Schema: &openapi3.SchemaRef{
						Value: typeToSchema[p.paramType],
					},
				},
			}
		})
	})
}

//func testHandler() handler {
//	mock := mockHandler{
//		request:  nil,
//		response: nil,
//	}
//	handler{
//		handlerFunc:    mock,
//		request:        mock.requestTypes(),
//		responses:      mock.response,
//		sourcePosition: sourcePosition{},
//	}
//	return
//}

const requestBody = "requestBody"
const pathParams = "pathParams"
const queryParams = "queryParams"

type mockHandler struct {
	request  map[string]reflect.Type
	response handlerResponses
}

func (t mockHandler) requestTypes() requestTypes {
	return requestTypes{
		requestBody: t.request[requestBody],
		pathParams:  t.request[pathParams],
		queryParams: t.request[queryParams]}
}
func (t mockHandler) responseTypes() handlerResponses {
	return t.response
}
func (t mockHandler) sourcePosition() sourcePosition {
	return sourcePosition{ok: true, file: "validations_test.go", line: 0}
}
func (t mockHandler) handlerFactory(openapi, BoundHandlerFunc) BoundHandlerFunc {
	return func(_ Context) (RawResponse, error) { return RawResponse{}, nil }
}
