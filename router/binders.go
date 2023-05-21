package router

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gin-gonic/gin/binding"
	"github.com/julienschmidt/httprouter"

	"github.com/piiano/cellotape/router/utils"
)

const contentTypeHeader = "Content-Type"

type binder[T any] func(*Context, *T) error

func nilBinder[T any](*Context, *T) error {
	return nil
}

// A request binder takes a Context with its untyped Context.Request and Context.Params and produce a typed Request.
type requestBinder[B, P, Q any] func(*Context) (Request[B, P, Q], error)

// A response binder takes a Context with its Context.Writer and previous Context.RawResponse to write a typed Response output.
type responseBinder[R any] func(*Context, Response[R]) (RawResponse, error)

// produce the binder function that can be called at runtime to create the httpRequest object for the handler.
func requestBinderFactory[B, P, Q any](oa openapi, types requestTypes) requestBinder[B, P, Q] {
	requestBodyBinder := requestBodyBinderFactory[B](types.requestBody, oa.contentTypes)
	pathParamsBinder := pathBinderFactory[P](types.pathParams)
	queryParamsBinder := queryBinderFactory[Q](types.queryParams)

	// this is what actually build the httpRequest object at runtime for the handler.
	return func(ctx *Context) (Request[B, P, Q], error) {
		var request = Request[B, P, Q]{
			Headers: ctx.Request.Header,
		}
		if err := requestBodyBinder(ctx, &request.Body); err != nil {
			return request, newBadRequestErr(ctx, err, InBody)
		}
		if err := pathParamsBinder(ctx, &request.PathParams); err != nil {
			return request, newBadRequestErr(ctx, err, InPathParams)
		}
		if err := queryParamsBinder(ctx, &request.QueryParams); err != nil {
			return request, newBadRequestErr(ctx, err, InQueryParams)
		}
		return request, nil
	}
}

// produce the httpRequest Body binder that can be used in runtime
func requestBodyBinderFactory[B any](requestBodyType reflect.Type, contentTypes ContentTypes) binder[B] {
	if requestBodyType == utils.NilType {
		return nilBinder[B]
	}
	return func(ctx *Context, body *B) error {
		input, err := validateBodyAndPopulateDefaults(ctx)
		if err != nil {
			return err
		}

		contentType, err := requestContentType(input.Request, contentTypes, JSONContentType{})
		if err != nil {
			return err
		}
		defer func() { _ = input.Request.Body.Close() }()
		bodyBytes, err := io.ReadAll(input.Request.Body)
		if err != nil {
			return err
		}
		if err = contentType.Decode(bodyBytes, body); err != nil {
			return err
		}
		return nil
	}
}

// validateBodyAndPopulateDefaults validate the request body with the openapi spec and populate the default values.
func validateBodyAndPopulateDefaults(ctx *Context) (*openapi3filter.RequestValidationInput, error) {
	input := requestValidationInput(ctx)
	if ctx.Operation.RequestBody != nil {
		if err := openapi3filter.ValidateRequestBody(ctx.Request.Context(), input, ctx.Operation.RequestBody.Value); err != nil {
			return nil, err
		}
	}
	return input, nil
}

// produce the pathParamInValue pathParams binder that can be used in runtime
func pathBinderFactory[P any](pathParamsType reflect.Type) binder[P] {
	if pathParamsType == utils.NilType {
		return nilBinder[P]
	}
	return func(ctx *Context, target *P) error {
		defaults, err := validateParamsAndPopulateDefaults(ctx, "path")
		if err != nil {
			return err
		}

		m := make(map[string][]string)
		for k, v := range defaults.PathParams {
			m[k] = []string{v}
		}

		if err = binding.Uri.BindUri(m, target); err != nil {
			return err
		}
		return nil
	}
}

// produce the queryParamInValue pathParams binder that can be used in runtime
func queryBinderFactory[Q any](queryParamsType reflect.Type) binder[Q] {
	if queryParamsType == utils.NilType {
		return nilBinder[Q]
	}
	paramFields := utils.StructKeys(queryParamsType, "form")
	nonArrayParams := utils.NewSet[string]()
	for param, paramType := range paramFields {
		kind := paramType.Type.Kind()
		if kind == reflect.Slice ||
			kind == reflect.Array ||
			(kind == reflect.Pointer &&
				(paramType.Type.Elem().Kind() == reflect.Slice ||
					paramType.Type.Elem().Kind() == reflect.Array)) {
			continue
		}
		nonArrayParams.Add(param)
	}

	return func(ctx *Context, queryParams *Q) error {
		defaults, err := validateParamsAndPopulateDefaults(ctx, "query")
		if err != nil {
			return err
		}

		if err = binding.Query.Bind(defaults.Request, queryParams); err != nil {
			return err
		}

		for param, values := range defaults.QueryParams {
			if nonArrayParams.Has(param) && len(values) > 1 {
				return fmt.Errorf("multiple values received for query param %s", param)
			}
		}
		return nil
	}
}

// responseBinderFactory creates a responseBinder that can be used in runtime
func responseBinderFactory[R any](responses handlerResponses, contentTypes ContentTypes) responseBinder[R] {
	return func(ctx *Context, r Response[R]) (RawResponse, error) {
		if ctx.RawResponse.Status != 0 {
			return *ctx.RawResponse, nil
		}
		contentType, err := responseContentType(r.contentType, contentTypes, JSONContentType{})
		if err != nil {
			log.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
		}
		responseType, exist := responses[r.status]
		if !exist {
			return RawResponse{}, fmt.Errorf("%w: %d", UnsupportedResponseStatusErr, r.status)
		}
		var responseBytes []byte
		if !responseType.isNilType {
			// this reflection call can not be avoided. we need some way to define multiple response types per handler
			// and struct fields is the only way to achieve that.
			responseField := reflect.ValueOf(r.response).FieldByIndex(responseType.fieldIndex).Interface()
			responseBytes, err = contentType.Encode(responseField)
			if err != nil {
				return RawResponse{}, err
			}
			r.headers.Set(contentTypeHeader, contentType.Mime())
		}
		bindResponseHeaders(ctx.Writer, r)
		ctx.Writer.WriteHeader(r.status)
		ctx.RawResponse.Status = r.status
		ctx.RawResponse.ContentType = r.contentType
		ctx.RawResponse.Body = responseBytes
		ctx.RawResponse.Headers = r.headers

		if _, err = ctx.Writer.Write(responseBytes); err != nil {
			return *ctx.RawResponse, err
		}

		validateResponse(ctx, r, responseBytes)

		return *ctx.RawResponse, nil
	}
}

// validateResponse validates the response against the spec. It logs a warning if the response violates the spec.
func validateResponse[R any](ctx *Context, r Response[R], responseBytes []byte) {
	input := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput(ctx),
		Status:                 r.status,
		Header:                 r.headers,
		Body:                   io.NopCloser(bytes.NewReader(responseBytes)),
		Options:                validationOptions(),
	}

	if err := openapi3filter.ValidateResponse(ctx.Request.Context(), input); err != nil {
		log.Printf("[WARNING] %s. response violates the spec\n", err)
	}
}

// bindResponseHeaders copies Response.headers to http.ResponseWriter.Header
func bindResponseHeaders[R any](writer http.ResponseWriter, r Response[R]) {
	for header, values := range r.headers {
		for _, value := range values {
			writer.Header().Add(header, value)
		}
	}
}

// requestContentType extracts the ContentType implementation to use base on the "Content-Type" request header.
// If "Content-Type" request header is missing fallback to the provided default ContentType.
func requestContentType(r *http.Request, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	mimeType := r.Header.Get(contentTypeHeader)

	if mimeType == "" {
		return defaultContentType, nil
	}

	parsedMimeType, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return nil, fmt.Errorf("%w: %q. %s", UnsupportedRequestContentTypeErr, mimeType, err)
	}

	if parsedMimeType == "*/*" {
		return defaultContentType, nil
	}

	if contentType, found := supportedTypes[parsedMimeType]; found {
		return contentType, nil
	}
	return nil, fmt.Errorf("%w: %q", UnsupportedRequestContentTypeErr, parsedMimeType)
}

// responseContentType extracts the ContentType implementation to use base on the response content type, supported content types and default fallback.
func responseContentType(responseContentType string, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	if responseContentType == "" {
		return defaultContentType, nil
	}
	if contentTypes, found := supportedTypes[responseContentType]; found {
		return contentTypes, nil
	}
	return defaultContentType, fmt.Errorf("%w: %s", UnsupportedResponseContentTypeErr, responseContentType)
}

func validateParamsAndPopulateDefaults(ctx *Context, in string) (*openapi3filter.RequestValidationInput, error) {
	input := requestValidationInput(ctx)
	parameters := utils.Filter(utils.Map(ctx.Operation.Parameters, func(p *openapi3.ParameterRef) *openapi3.Parameter {
		return p.Value
	}), func(p *openapi3.Parameter) bool { return p.In == in })

	for _, param := range parameters {
		if err := openapi3filter.ValidateParameter(ctx.Request.Context(), input, param); err != nil {
			return nil, err
		}
	}

	// after processing params input is populated with defaults
	return input, nil
}

func requestValidationInput(ctx *Context) *openapi3filter.RequestValidationInput {
	input := openapi3filter.RequestValidationInput{
		Request:     &http.Request{},
		PathParams:  make(map[string]string),
		QueryParams: url.Values{},
		Options:     validationOptions(),
		Route: &routers.Route{
			Operation: ctx.Operation.Operation,
		},
		ParamDecoder: nil,
	}

	if ctx.Request != nil {
		input.Request = ctx.Request
		if ctx.Request.URL != nil {
			input.QueryParams = ctx.Request.URL.Query()
		}
	}

	if ctx.Params != nil {
		input.PathParams = utils.FromEntries(utils.Map(*ctx.Params, func(p httprouter.Param) utils.Entry[string, string] {
			return utils.Entry[string, string]{
				Key:   p.Key,
				Value: p.Value,
			}
		}))
	}

	return &input
}

func validationOptions() *openapi3filter.Options {
	options := openapi3filter.Options{}
	// Customize the error message returned by the kin-openapi library to be more user-friendly.
	options.WithCustomSchemaErrorFunc(func(err *openapi3.SchemaError) string {
		return err.Reason
	})
	return &options
}
