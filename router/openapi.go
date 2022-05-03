package router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/piiano/restcontroller/schema_validator"
	"github.com/piiano/restcontroller/utils"
	"log"
	"net/http"
	"regexp"
)

type OpenAPIRouter interface {
	// Use middlewares on the root handler
	Use(...Handler) OpenAPIRouter
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) OpenAPIRouter
	// WithOperation attaches an Handler for an operation ID string on the OpenAPISpec
	WithOperation(string, Handler, ...Handler) OpenAPIRouter
	// WithContentType Add support for encoding additional content type
	WithContentType(ContentType) OpenAPIRouter
	// AsHandler validates the validity of the specified implementation with the registered OpenAPISpec
	// returns a http.Handler that implements nil value if all checks passed correctly
	AsHandler() (http.Handler, error)
}

// NewOpenAPIRouter creates new OpenAPIRouter router with the provided OpenAPISpec
func NewOpenAPIRouter(spec OpenAPISpec) OpenAPIRouter {
	return NewOpenAPIRouterWithOptions(spec, Options{
		RecoverOnPanic: true,
	})
}

// NewOpenAPIRouterWithOptions creates new OpenAPIRouter router with the provided OpenAPISpec and custom options
func NewOpenAPIRouterWithOptions(spec OpenAPISpec, options Options) OpenAPIRouter {
	return &openapi{
		spec:         spec,
		options:      options,
		contentTypes: DefaultContentTypes(),
		group: group{
			handlers:   []Handler{},
			groups:     []group{},
			operations: map[string]operation{},
		},
	}
}

type Options struct {
	// TODO: Add support for options
	// Fine tune schema validation during initialization.
	InitializationSchemaValidation schema_validator.Options

	// provide a default recover from panic that return status 500
	RecoverOnPanic bool
}

type openapi struct {
	spec         OpenAPISpec
	options      Options
	contentTypes ContentTypes
	group
}

func (oa *openapi) Use(handlers ...Handler) OpenAPIRouter {
	oa.group.Use(handlers...)
	return oa
}
func (oa *openapi) WithGroup(group Group) OpenAPIRouter {
	oa.group.WithGroup(group)
	return oa
}
func (oa *openapi) WithOperation(id string, handlerFunc Handler, handlers ...Handler) OpenAPIRouter {
	oa.group.WithOperation(id, handlerFunc, handlers...)
	return oa
}
func (oa *openapi) WithContentType(contentType ContentType) OpenAPIRouter {
	oa.contentTypes[contentType.Mime()] = contentType
	return oa
}

func (oa *openapi) AsHandler() (http.Handler, error) {
	pathParamsMatcher := regexp.MustCompile(`\{([^/]*)}`)
	errs := utils.NewErrorsCollector()
	router := httprouter.New()
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
		chainResponseTypes := utils.Map(append(flatOp.handlers, flatOp.operationHandler), func(handler Handler) handlerResponseTypes {
			return handler.responseTypes()
		})
		errs.AddIfNotNil(flatOp.validateOperationTypes(*oa, chainResponseTypes))
		if errs.ErrorOrNil() != nil {
			continue
		}
		path := pathParamsMatcher.ReplaceAllString(specOp.path, ":$1")
		handlersChainHead := createHandlersChain(*oa, append(flatOp.handlers, flatOp.operationHandler)...)
		router.Handle(specOp.method, path, handlersChainHead)
	}
	for _, pathItem := range oa.spec.Paths {
		for _, specOp := range pathItem.Operations() {
			if !declaredOperation[specOp.OperationID] {
				errs.AddIfNotNil(fmt.Errorf("missing handler for operation id %q", specOp.OperationID))
			}
		}
	}
	return router, errs.ErrorOrNil()
}

func createHandlersChain(oa openapi, handlers ...Handler) (head httprouter.Handle) {
	next := func(HandlerContext) (Response[any], error) { return Response[any]{}, nil }
	for i := len(handlers) - 1; i >= 0; i-- {
		handler := handlers[i]
		next = handler.handler(oa, next)
	}
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		if oa.options.RecoverOnPanic {
			defer func() {
				if r := recover(); r != nil {
					writer.WriteHeader(500)
					log.Printf("[ERROR] recovered from panic. %v. respond with status 500\n", r)
				}
			}()
		}
		response, err := next(HandlerContext{Request: request, Params: params})
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		for header, values := range response.Headers {
			for _, value := range values {
				writer.Header().Add(header, value)
			}
		}
		writer.WriteHeader(response.Status)
		writer.Write(response.Bytes)
	}
}

// flattenOperations takes a group with separate operations, handlers, and nested groups and flatten them into a flat
// Operation slice that include for each Operation its own NewOperationHandler, attached handlers, and attached group handlers.
func flattenOperations(g group) []operation {
	flatOperations := make([]operation, 0)
	for id, op := range g.operations {
		flatOperations = append(flatOperations, operation{
			id:               id,
			handlers:         utils.ConcatSlices(g.handlers, op.handlers),
			operationHandler: op.operationHandler,
			responseTypes:    op.responseTypes,
		})
	}
	for _, nestedGroup := range g.groups {
		for _, flatOp := range flattenOperations(nestedGroup) {
			flatOperations = append(flatOperations, operation{
				id:               flatOp.id,
				handlers:         utils.ConcatSlices(g.handlers, flatOp.handlers),
				operationHandler: flatOp.operationHandler,
				responseTypes:    flatOp.responseTypes,
			})
		}
	}
	return flatOperations
}
