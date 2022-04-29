package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/piiano/restcontroller/utils"
	"net/http"
	"regexp"
)

type OpenAPI interface {
	// Use middlewares on the root handler
	Use(...http.HandlerFunc) OpenAPI
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) OpenAPI
	// WithOperation attaches an OperationHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, OperationHandler, ...http.HandlerFunc) OpenAPI
	// WithResponse Define a reusable response
	// Responses are validated with the OpenAPISpec when calling Validate with proper options
	WithResponse(HttpResponse) OpenAPI
	// WithContentType Add support for encoding additional content type
	WithContentType(ContentType) OpenAPI
	// AsHandler validates the validity of the specified implementation with the registered OpenAPISpec
	// returns a http.Handler that implements nil value if all checks passed correctly
	AsHandler() (http.HandlerFunc, error)

	getSpec() OpenAPISpec
	getOptions() Options
	getContentTypes() ContentTypes
	groupInternals
}

// NewOpenAPI creates new OpenAPI router with the provided OpenAPISpec
func NewOpenAPI(spec OpenAPISpec) OpenAPI {
	return &openapi{
		spec:         spec,
		options:      Options{},
		handlers:     []http.HandlerFunc{},
		groups:       []groupInternals{},
		operations:   map[string]Operation{},
		contentTypes: ContentTypes{},
		responses:    []HttpResponse{},
	}
}

// NewOpenAPIWithOptions creates new OpenAPI router with the provided OpenAPISpec and custom options
func NewOpenAPIWithOptions(spec OpenAPISpec, options Options) OpenAPI {
	return &openapi{
		spec:         spec,
		options:      options,
		handlers:     []http.HandlerFunc{},
		groups:       []groupInternals{},
		operations:   map[string]Operation{},
		contentTypes: ContentTypes{},
		responses:    []HttpResponse{},
	}
}

type openapi struct {
	spec         OpenAPISpec
	options      Options
	handlers     []http.HandlerFunc
	groups       []groupInternals
	operations   map[string]Operation
	contentTypes ContentTypes
	responses    []HttpResponse
}

func (oa *openapi) Use(handlers ...http.HandlerFunc) OpenAPI {
	oa.handlers = append(oa.handlers, handlers...)
	return oa
}
func (oa *openapi) WithGroup(group Group) OpenAPI {
	oa.groups = append(oa.groups, group)
	return oa
}
func (oa *openapi) WithOperation(id string, handlerFunc OperationHandler, handlers ...http.HandlerFunc) OpenAPI {
	oa.operations[id] = operation{
		operationHandler: handlerFunc,
		handlers:         handlers,
	}
	return oa
}
func (oa *openapi) WithResponse(response HttpResponse) OpenAPI {
	oa.responses = append(oa.responses, response)
	return oa
}
func (oa *openapi) WithContentType(contentType ContentType) OpenAPI {
	oa.contentTypes[contentType.Mime()] = contentType
	return oa
}

func (oa *openapi) AsHandler() (http.HandlerFunc, error) {
	pathParamsMatcher := regexp.MustCompile(`\{([^/]*)}`)
	errs := utils.NewErrorsCollector()
	engine := gin.New()
	flatOperations := flattenOperations(oa)
	declaredOperation := make(map[string]bool, len(flatOperations))
	for _, flatOp := range flatOperations {
		if declaredOperation[flatOp.getId()] {
			errs.AddIfNotNil(fmt.Errorf("multiple handlers found for operation id %q", flatOp.getId()))
			continue
		}
		declaredOperation[flatOp.getId()] = true
		specOp, found := oa.getSpec().findSpecOperationByID(flatOp.getId())
		if !found {
			errs.AddIfNotNil(fmt.Errorf("handler recieved for non exising operation id %q is spec", flatOp.getId()))
		}
		errs.AddIfNotNil(flatOp.validateHandlerTypes(oa))
		handler := flatOp.getHandlerFunc().asGinHandler(oa)
		path := pathParamsMatcher.ReplaceAllString(specOp.path, ":$1")
		engine.Handle(specOp.method, path, append(utils.Map(flatOp.getHandlers(), func(h http.HandlerFunc) gin.HandlerFunc {
			return gin.WrapH(h)
		}), handler)...)
	}
	for _, pathItem := range oa.getSpec().Paths {
		for _, specOp := range pathItem.Operations() {
			if !declaredOperation[specOp.OperationID] {
				errs.AddIfNotNil(fmt.Errorf("missing handler for operation id %q", specOp.OperationID))
			}
		}
	}
	return engine.ServeHTTP, errs.ErrorOrNil()
}

func (oa openapi) getGroups() []groupInternals         { return oa.groups }
func (oa openapi) getHandlers() []http.HandlerFunc     { return oa.handlers }
func (oa openapi) getOperations() map[string]Operation { return oa.operations }
func (oa *openapi) getSpec() OpenAPISpec               { return oa.spec }
func (oa *openapi) getOptions() Options                { return oa.options }
func (oa *openapi) getContentTypes() ContentTypes      { return oa.contentTypes }

//specResponses := oa.spec.Components.Responses
////if oa.options.FailOnMissing.Responses && len(specResponses) != len(oa.responses) {
////	return nil, fmt.Errorf("spec responses (%d) don't match provided responses (%d)", len(specResponses), len(oa.responses))
////}
//for _, r := range oa.responses {
//	respRef := specResponses.Get(r.getStatus())
//	if respRef == nil {
//		return nil, fmt.Errorf("response for status %d is missing in the spec", r.getStatus())
//	}
//	contentType := r.getContentType()
//	respMedia := respRef.Value.Content.Get(contentType)
//	if respMedia == nil {
//		return nil, fmt.Errorf("response for status %d is missing content type %q in the spec", r.getStatus(), contentType)
//	}
//	r.getBodyType()
//	openapi3.NewSchema()
//}
