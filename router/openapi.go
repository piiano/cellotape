package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/piiano/restcontroller/utils"
	"net/http"
	"regexp"
)

type OpenAPIRouter interface {
	// Use middlewares on the root handler
	Use(...http.HandlerFunc) OpenAPIRouter
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) OpenAPIRouter
	// WithOperation attaches an OperationHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, OperationHandler, ...http.HandlerFunc) OpenAPIRouter
	// WithContentType Add support for encoding additional content type
	WithContentType(ContentType) OpenAPIRouter
	// AsHandler validates the validity of the specified implementation with the registered OpenAPISpec
	// returns a http.Handler that implements nil value if all checks passed correctly
	AsHandler() (http.HandlerFunc, error)
}

// NewOpenAPIRouter creates new OpenAPIRouter router with the provided OpenAPISpec
func NewOpenAPIRouter(spec OpenAPISpec) OpenAPIRouter {
	return NewOpenAPIRouterWithOptions(spec, Options{})
}

// NewOpenAPIRouterWithOptions creates new OpenAPIRouter router with the provided OpenAPISpec and custom options
func NewOpenAPIRouterWithOptions(spec OpenAPISpec, options Options) OpenAPIRouter {
	return &openapi{
		spec:         spec,
		options:      options,
		contentTypes: DefaultContentTypes(),
		group: group{
			handlers:   []http.HandlerFunc{},
			groups:     []group{},
			operations: map[string]operation{},
		},
	}
}

type openapi struct {
	spec         OpenAPISpec
	options      Options
	contentTypes ContentTypes
	group
}

func (oa *openapi) Use(handlers ...http.HandlerFunc) OpenAPIRouter {
	oa.group.Use(handlers...)
	return oa
}
func (oa *openapi) WithGroup(group Group) OpenAPIRouter {
	oa.group.WithGroup(group)
	return oa
}
func (oa *openapi) WithOperation(id string, handlerFunc OperationHandler, handlers ...http.HandlerFunc) OpenAPIRouter {
	oa.group.WithOperation(id, handlerFunc, handlers...)
	return oa
}
func (oa *openapi) WithContentType(contentType ContentType) OpenAPIRouter {
	oa.contentTypes[contentType.Mime()] = contentType
	return oa
}

func (oa *openapi) AsHandler() (http.HandlerFunc, error) {
	pathParamsMatcher := regexp.MustCompile(`\{([^/]*)}`)
	errs := utils.NewErrorsCollector()
	engine := gin.New()
	flatOperations := flattenOperations(oa.group)
	declaredOperation := make(map[string]bool, len(flatOperations))
	for _, flatOp := range flatOperations {
		if declaredOperation[flatOp.id] {
			errs.AddIfNotNil(fmt.Errorf("multiple handlers found for operation id %q", flatOp.id))
			continue
		}
		declaredOperation[flatOp.id] = true
		specOp, found := oa.spec.findSpecOperationByID(flatOp.id)
		if !found {
			errs.AddIfNotNil(fmt.Errorf("handler recieved for non exising operation id %q is spec", flatOp.id))
		}
		errs.AddIfNotNil(flatOp.validateHandlerTypes(*oa))
		if errs.ErrorOrNil() != nil {
			continue
		}
		handler := flatOp.operationHandler.asGinHandler(*oa)
		path := pathParamsMatcher.ReplaceAllString(specOp.path, ":$1")
		engine.Handle(specOp.method, path, append(utils.Map(flatOp.handlers, func(h http.HandlerFunc) gin.HandlerFunc {
			return gin.WrapH(h)
		}), handler)...)
	}
	for _, pathItem := range oa.spec.Paths {
		for _, specOp := range pathItem.Operations() {
			if !declaredOperation[specOp.OperationID] {
				errs.AddIfNotNil(fmt.Errorf("missing handler for operation id %q", specOp.OperationID))
			}
		}
	}
	return engine.ServeHTTP, errs.ErrorOrNil()
}

// flattenOperations takes a group with separate operations, handlers, and nested groups and flatten them into a flat
// Operation slice that include for each Operation its own OperationFunc, attached handlers, and attached group handlers.
func flattenOperations(g group) []operation {
	flatOperations := make([]operation, 0)
	for id, op := range g.operations {
		flatOperations = append(flatOperations, operation{
			id:               id,
			handlers:         utils.ConcatSlices(g.handlers, op.handlers),
			operationHandler: op.operationHandler,
		})
	}
	for _, nestedGroup := range g.groups {
		for _, flatOp := range flattenOperations(nestedGroup) {
			flatOperations = append(flatOperations, operation{
				id:               flatOp.id,
				handlers:         utils.ConcatSlices(g.handlers, flatOp.handlers),
				operationHandler: flatOp.operationHandler,
			})
		}
	}
	return flatOperations
}
