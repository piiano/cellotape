package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"net/http"
)

type OpenAPI interface {
	// WithContext Adds a go context that is passed to OperationHandler
	WithContext(context.Context) OpenAPI
	// WithOptions Sets options for the openapi validation and routing mechanism
	WithOptions(Options) OpenAPI
	// Use middlewares on the root handler
	Use(...http.Handler) OpenAPI
	// WithGroup defines a virtual group that allow defining middlewares for group of routes more easily
	WithGroup(Group) OpenAPI
	// WithOperation attaches an OperationHandler for an operation ID string on the OpenAPISpec
	WithOperation(string, OperationHandler, ...http.Handler) OpenAPI
	// WithResponse Define a reusable response
	// Responses are validated with the OpenAPISpec when calling Validate with proper options
	WithResponse(HttpResponse) OpenAPI
	// WithContentType Add support for encoding additional content type
	WithContentType(ContentType) OpenAPI
	// WithSpec Sets the OpenAPISpec to be used for OpenAPI router and validations
	WithSpec(OpenAPISpec) OpenAPI
	// Spec returns the OpenAPISpec object
	Spec() OpenAPISpec
	// AsHandler validates the validity of the specified implementation with the registered OpenAPISpec
	// returns a http.Handler that implements nil value if all checks passed correctly
	AsHandler() (http.Handler, error)
}

func NewOpenAPI() OpenAPI {
	return openapi{
		spec:       OpenAPISpec{},
		operations: map[string]operationHandler{},
		handlers:   []http.Handler{},
		ctx:        nil,
		responses:  []HttpResponse{},
		options:    Options{},
		errors:     []error{},
	}
}

type openapi struct {
	spec         OpenAPISpec
	errors       []error
	options      Options
	ctx          context.Context
	operations   map[string]operationHandler
	contentTypes map[string]ContentType
	handlers     []http.Handler
	responses    []HttpResponse
}

func (oa openapi) WithContentType(contentType ContentType) OpenAPI {
	oa.contentTypes[contentType.Mime()] = contentType
	return oa
}

type operationHandler struct {
	operationHandler OperationHandler
	handlers         []http.Handler
}

func (oa *openapi) appendErrIfNotNil(err error) {
	if err != nil {
		oa.errors = append(oa.errors, err)
	}
}

func (oa openapi) WithContext(ctx context.Context) OpenAPI {
	oa.ctx = ctx
	return oa
}
func (oa openapi) WithOptions(options Options) OpenAPI {
	oa.options = options
	return oa
}
func (oa openapi) Use(handlers ...http.Handler) OpenAPI {
	oa.handlers = append(oa.handlers, handlers...)
	return oa
}
func (oa openapi) WithGroup(group Group) OpenAPI {
	for _, err := range group.getErrors() {
		oa.appendErrIfNotNil(err)
	}
	handlers := group.getHandlers()
	for id, oh := range group.getOperations() {
		if _, ok := oa.operations[id]; ok {
			oa.appendErrIfNotNil(fmt.Errorf("operation with id %q already declared", id))
			continue
		}
		oa.operations[id] = operationHandler{
			handlers:         append(handlers, oh.handlers...),
			operationHandler: oh.operationHandler,
		}
	}
	return oa
}
func (oa openapi) WithOperation(id string, handler OperationHandler, handlers ...http.Handler) OpenAPI {
	oa.operations[id] = operationHandler{
		operationHandler: handler,
		handlers:         handlers,
	}
	return oa
}
func (oa openapi) WithResponse(response HttpResponse) OpenAPI {
	oa.responses = append(oa.responses, response)
	return oa
}
func (oa openapi) WithSpec(spec OpenAPISpec) OpenAPI {
	bytes, err := json.Marshal(spec)
	oa.appendErrIfNotNil(err)
	clonedSpec, err := openapi3.NewLoader().LoadFromData(bytes)
	oa.appendErrIfNotNil(err)
	oa.spec = OpenAPISpec(*clonedSpec)
	return oa
}
func (oa openapi) Spec() OpenAPISpec {
	return oa.spec
}

// TODO: implement
func (oa openapi) AsHandler() (http.Handler, error) {
	if len(oa.errors) > 0 {
		return nil, oa.errors[0]
	}
	specResponses := oa.spec.Components.Responses
	if oa.options.FailOnMissing.Responses && len(specResponses) != len(oa.responses) {
		return nil, fmt.Errorf("spec responses (%d) don't match provided responses (%d)", len(specResponses), len(oa.responses))
	}
	for _, r := range oa.responses {
		respRef := specResponses.Get(r.getStatus())
		if respRef == nil {
			return nil, fmt.Errorf("response for status %d is missing in the spec", r.getStatus())
		}
		contentType := r.getContentType()
		respMedia := respRef.Value.Content.Get(contentType)
		if respMedia == nil {
			return nil, fmt.Errorf("response for status %d is missing content type %q in the spec", r.getStatus(), contentType)
		}
		r.getBodyType()
		openapi3.NewSchema()
	}
	for id, oh := range oa.operations {
		if _, ok := oa.operations[id]; ok {
			oa.appendErrIfNotNil(fmt.Errorf("operation with id %q already declared", id))
			continue
		}
		oa.operations[id] = operationHandler{
			handlers:         append(oa.handlers, oh.handlers...),
			operationHandler: oh.operationHandler,
		}
	}

	panic(errors.New("unimplemented"))
}
