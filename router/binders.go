package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/julienschmidt/httprouter"

	"github.com/piiano/cellotape/router/utils"
)

const contentTypeHeader = "Content-Type"

// A request binder takes a Context with its untyped Context.Request and Context.Params and produce a typed Request.
type requestBinder[B, P, Q any] func(ctx Context) (Request[B, P, Q], error)

// A response binder takes a Context with its Context.Writer and previous Context.RawResponse to write a typed Response output.
type responseBinder[R any] func(Context, Response[R]) (RawResponse, error)

// produce the binder function that can be called at runtime to create the httpRequest object for the handler.
func requestBinderFactory[B, P, Q any](oa openapi, types requestTypes) requestBinder[B, P, Q] {
	requestBodyBinder := requestBodyBinderFactory[B](types.requestBody, oa.contentTypes)
	pathParamsBinder := pathBinderFactory[P](types.pathParams)
	queryParamsBinder := queryBinderFactory[Q](types.queryParams)

	// this is what actually build the httpRequest object at runtime for the handler.
	return func(ctx Context) (Request[B, P, Q], error) {
		var request = Request[B, P, Q]{
			Headers: ctx.Request.Header,
		}
		if err := requestBodyBinder(ctx.Request, &request.Body); err != nil {
			return request, newBadRequestErr(ctx, err, InBody)
		}
		if err := pathParamsBinder(ctx.Params, &request.PathParams); err != nil {
			return request, newBadRequestErr(ctx, err, InPathParams)
		}
		if err := queryParamsBinder(ctx.Request, &request.QueryParams); err != nil {
			return request, newBadRequestErr(ctx, err, InQueryParams)
		}
		return request, nil
	}
}

// produce the httpRequest Body binder that can be used in runtime
func requestBodyBinderFactory[B any](requestBodyType reflect.Type, contentTypes ContentTypes) func(*http.Request, *B) error {
	if requestBodyType == nilType {
		return func(r *http.Request, body *B) error { return nil }
	}
	return func(r *http.Request, body *B) error {
		contentType, err := requestContentType(r, contentTypes, JSONContentType{})
		if err != nil {
			return err
		}
		defer func() { _ = r.Body.Close() }()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		if err = contentType.Decode(bodyBytes, body); err != nil {
			return err
		}
		return nil
	}
}

// produce the pathParamInValue pathParams binder that can be used in runtime
func pathBinderFactory[P any](pathParamsType reflect.Type) func(*httprouter.Params, *P) error {
	if pathParamsType == nilType {
		return func(params *httprouter.Params, body *P) error { return nil }
	}
	return func(params *httprouter.Params, pathParams *P) error {
		m := make(map[string][]string)
		for _, v := range *params {
			m[v.Key] = []string{v.Value}
		}
		if err := binding.Uri.BindUri(m, pathParams); err != nil {
			return err
		}
		return nil
	}
}

// produce the queryParamInValue pathParams binder that can be used in runtime
func queryBinderFactory[Q any](queryParamsType reflect.Type) func(*http.Request, *Q) error {
	if queryParamsType == nilType {
		return func(*http.Request, *Q) error { return nil }
	}
	paramFields := structKeys(queryParamsType, "form")
	nonArrayParams := utils.NewSet[string]()
	for param, paramType := range paramFields {
		if paramType.Type.Kind() == reflect.Slice ||
			paramType.Type.Kind() == reflect.Array ||
			(paramType.Type.Kind() == reflect.Pointer &&
				(paramType.Type.Elem().Kind() == reflect.Slice ||
					paramType.Type.Elem().Kind() == reflect.Array)) {
			continue
		}
		nonArrayParams.Add(param)
	}

	return func(r *http.Request, queryParams *Q) error {
		if err := binding.Query.Bind(r, queryParams); err != nil {
			return err
		}
		for param, values := range r.URL.Query() {
			if nonArrayParams.Has(param) && len(values) > 1 {
				return fmt.Errorf("multiple values received for query param %s", param)
			}
		}
		return nil
	}
}

// responseBinderFactory creates a responseBinder that can be used in runtime
func responseBinderFactory[R any](responses handlerResponses, contentTypes ContentTypes) responseBinder[R] {
	return func(ctx Context, r Response[R]) (RawResponse, error) {
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
		_, err = ctx.Writer.Write(responseBytes)
		ctx.RawResponse.Status = r.status
		ctx.RawResponse.ContentType = r.contentType
		ctx.RawResponse.Body = responseBytes
		ctx.RawResponse.Headers = r.headers
		return *ctx.RawResponse, err
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
	if mimeType == "*/*" || mimeType == "" {
		return defaultContentType, nil
	}
	if contentType, found := supportedTypes[mimeType]; found {
		return contentType, nil
	}
	return nil, fmt.Errorf("%w: %q", UnsupportedRequestContentTypeErr, mimeType)
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
