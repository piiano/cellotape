package router

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/julienschmidt/httprouter"

	"github.com/piiano/cellotape/router/utils"
)

func createMainRouterHandler(oa *openapi) (http.Handler, error) {
	flatOperations := flattenOperations(oa.group)
	if err := validateOpenAPIRouter(oa, flatOperations); err != nil {
		return nil, err
	}
	router := httprouter.New()
	router.HandleMethodNotAllowed = false

	logger := oa.logger()
	pathParamsMatcher := regexp.MustCompile(`\{([^/}]*)}`)

	specOperations := oa.spec.Operations()
	for _, flatOp := range flatOperations {
		specOp := specOperations[flatOp.id]
		path := pathParamsMatcher.ReplaceAllString(specOp.Path, ":$1")
		chainHead := chainHandlers(*oa, append(flatOp.handlers, flatOp.handler)...)
		httpRouterHandler := asHttpRouterHandler(*oa, specOp, chainHead)
		router.Handle(specOp.Method, path, httpRouterHandler)
		logger.Infof("register handler for operation %q - %s %s", flatOp.id, specOp.Method, specOp.Path)
	}

	setOptionsHandlers(router, oa)

	// For Kin-openapi to be able to validate a request and set default values it need to know how to decode and encode
	// the request body for any supported content type.
	for _, contentType := range oa.contentTypes {
		mimeType := contentType.Mime()
		if openapi3filter.RegisteredBodyEncoder(mimeType) == nil {
			openapi3filter.RegisterBodyEncoder(contentType.Mime(), contentType.Encode)
		}
		if openapi3filter.RegisteredBodyDecoder(mimeType) == nil {
			openapi3filter.RegisterBodyDecoder(contentType.Mime(), createDecoder(contentType))
		}
	}

	registerAdditionalOpenAPIFormatValidations()

	return router, nil
}

// setOptionsHandlers sets the OPTIONS handlers for the router.
func setOptionsHandlers(router *httprouter.Router, oa *openapi) {
	router.HandleOPTIONS = oa.options.OptionsHandler != nil

	// If HandleOptions is nil, NotFound simply return 404 for any request.
	// When HandleOptions is not nil, NotFound for any OPTIONS request will call the global OPTIONS handler to return the "Allow: OPTIONS" header.
	// For any other request it will return 404
	router.NotFound = notFoundHandler(oa)

	if !router.HandleOPTIONS {
		return
	}

	// GlobalOPTIONS is used for the OPTIONS handler for all paths.
	// It receives all registered methods for a path in the "Allow" header.
	router.GlobalOPTIONS = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		oa.options.OptionsHandler.ServeHTTP(writer, request)
	})
}

func createDecoder(contentType ContentType) func(reader io.Reader, _ http.Header, schema *openapi3.SchemaRef, enc openapi3filter.EncodingFn) (any, error) {
	return func(reader io.Reader, _ http.Header, schema *openapi3.SchemaRef, enc openapi3filter.EncodingFn) (any, error) {
		bytes, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		var target any
		if err = contentType.Decode(bytes, &target); err != nil {
			return nil, err
		}

		// For kin-openapi to be able to validate a request it requires that the decoded value will on of
		// the values received when decoding JSON to any.
		// e.g. any, []any, []map[string]any, etc.
		//
		// After using the custom decoder we get a value of the type of the target struct.
		// To overcome this we marshal the target to JSON and then unmarshal it to any.

		jsonBytes, err := json.Marshal(target)
		if err != nil {
			return nil, err
		}

		var jsonValue any
		if err = json.Unmarshal(jsonBytes, &jsonValue); err != nil {
			return nil, err
		}

		return jsonValue, nil
	}
}

func (oa *openapi) logger() utils.Logger {
	return utils.NewLoggerWithLevel(oa.options.LogOutput, oa.options.LogLevel)
}

// flattenOperations takes a group with separate operations, handlers, and nested groups and flatten them into a flat
// OperationHandler slice that include for each OperationHandler its own NewOperation, attached handlers, and attached group handlers.
func flattenOperations(g group) []operation {
	groupsOperations := utils.ConcatSlices[operation](utils.Map(g.groups, flattenOperations)...)
	return utils.Map(append(g.operations, groupsOperations...), func(op operation) operation {
		return operation{
			id:       op.id,
			handler:  op.handler,
			handlers: utils.ConcatSlices(g.handlers, op.handlers),
		}
	})
}

func chainHandlers(oa openapi, handlers ...handler) (head BoundHandlerFunc) {
	var next BoundHandlerFunc
	for i := len(handlers) - 1; i >= 0; i-- {
		next = handlers[i].handlerFunc.handlerFactory(oa, next)
	}
	next = ErrorHandler(func(c *Context, err error) (Response[any], error) {
		var badRequestError BadRequestErr
		if err != nil && c.RawResponse.Status == 0 && errors.As(err, &badRequestError) {
			c.Writer.Header().Add("Content-Type", "text/plain")
			c.Writer.WriteHeader(400)
			_, writeErr := c.Writer.Write([]byte(err.Error()))
			return Error[any](writeErr)
		}
		return Error[any](err)
	}).handlerFactory(oa, next)
	return next
}

func asHttpRouterHandler(oa openapi, specOp SpecOperation, head BoundHandlerFunc) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		if oa.options.RecoverOnPanic {
			defer defaultRecoverBehaviour(writer)
		}
		ctx := &Context{
			Operation:   specOp,
			Writer:      writer,
			Request:     request,
			Params:      &params,
			RawResponse: &RawResponse{Status: 0},
		}

		_, err := head(ctx)
		if err != nil || ctx.RawResponse.Status == 0 {
			writer.WriteHeader(500)
			return
		}
	}
}

func defaultRecoverBehaviour(writer http.ResponseWriter) {
	if r := recover(); r != nil {
		writer.WriteHeader(500)
		log.Printf("[Error] recovered from panic. %v. respond with status 500\n", r)
		debug.PrintStack()
	}
}

// DefaultOptionsHandler This handler will be called for any OPTIONS request with an "Allow" header that will include all the methods that are defined for the path.
func DefaultOptionsHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
}

func notFoundHandler(oa *openapi) http.HandlerFunc {
	if oa.options.OptionsHandler == nil {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodOptions {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		// If we get to NotFound it means that the path does not exist so we can set the "Allow" header to "OPTIONS" only and call the OPTIONS handler.
		writer.Header().Set("Allow", "OPTIONS")
		oa.options.OptionsHandler.ServeHTTP(writer, request)
	}
}
