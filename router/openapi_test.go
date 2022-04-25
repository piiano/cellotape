package router

import (
	_ "embed"
	"github.com/getkin/kin-openapi/openapi3"
	"reflect"
	"testing"
)

func TestNewOpenAPI(t *testing.T) {
	oa := NewOpenAPI()
	oa, ok := oa.(OpenAPI)
	if !ok {
		t.Error("NewOpenAPI returned value doesn't implement the OpenAPI interface")
	}
}
func TestNewOpenAPIDefaultSpec(t *testing.T) {
	oa := NewOpenAPI()
	spec := oa.Spec()
	emptySpec := OpenAPISpec{}
	if !reflect.DeepEqual(spec, emptySpec) {
		t.Error("Spec() should be initialized to empty spec")
	}
}
func TestNewOpenAPIWithSpec(t *testing.T) {
	openapi := NewOpenAPI()
	testTitle := "Test API"
	testSpec := OpenAPISpec{
		Info: &openapi3.Info{Title: testTitle},
	}
	newOpenapi := openapi.WithSpec(testSpec)
	spec := newOpenapi.Spec()
	if &newOpenapi == &openapi {
		t.Error("WithSpec() should be immutable")
	}
	if &spec == &testSpec {
		t.Error("Spec() should return a copy of the spec provided to WithSpec()")
	}
	if spec.Info.Title != testTitle {
		t.Error("Spec() should hold a spec with value provided to WithSpec()")
	}
}

//
////go:embed hello-world-openapi.yaml
//var Spec []byte

//func TestExampleMainFunction(f *testing.T) {
//	testSpec := OpenAPISpec{
//		Info: &openapi3.Info{Title: "testTitle"},
//	}
//	oa := NewOpenAPI().
//		WithSpec(testSpec).
//		Use().
//		//WithGroup(NewGroup().
//		//Use().
//		//WithOperation("greet", example.GreetOperationHandler).
//		//WithOperation("greet", example.GreetOperationHandler).
//		//WithOperation("greet", example.GreetOperationHandler),
//		//).
//		//WithContentType().
//		WithResponse(NewHttpResponse[error](400, "application/json")).
//		WithResponse(NewHttpResponse[error](500, "application/json")).
//		WithResponse(NewHttpResponse[error](404, "application/json"))
//
//	handler, err := oa.AsHandler()
//	if err != nil {
//		panic(err)
//	}
//	if err := http.ListenAndServe(":8080", handler); err != nil {
//		panic(err)
//	}
//}
