package router

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"reflect"
)

const ContentTypeHeader = "Content-Type"
const AcceptHeader = "Accept"

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
			Context: ctx.Request.Context(),
			Method:  ctx.Request.Method,
			URL:     ctx.Request.URL,
			Headers: ctx.Request.Header,
		}
		if err := requestBodyBinder(ctx.Request, &request.Body); err != nil {
			return request, err
		}
		if err := pathParamsBinder(ctx.Params, &request.PathParams); err != nil {
			return request, err
		}
		if err := queryParamsBinder(ctx.Request, &request.QueryParams); err != nil {
			return request, err
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
		contentType, err := requestContentType(r, contentTypes, JsonContentType{})
		if err != nil {
			return err
		}
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
		return binding.Uri.BindUri(m, pathParams)
	}
}

// produce the queryParamInValue pathParams binder that can be used in runtime
func queryBinderFactory[Q any](queryParamsType reflect.Type) func(*http.Request, *Q) error {
	if queryParamsType == nilType {
		return func(*http.Request, *Q) error { return nil }
	}
	return func(r *http.Request, queryParams *Q) error {
		if err := binding.Query.Bind(r, queryParams); err != nil {
			return err
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
		contentType, err := responseContentType(ctx.Request, contentTypes, JsonContentType{})
		if err != nil {
			log.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
		}
		responseType, exist := responses[r.Status]
		if !exist {
			return RawResponse{}, fmt.Errorf("status %d is not part of the possible operation responses", r.Status)
		}
		var responseBytes []byte
		if !responseType.isNilType {
			// this reflection call can not be avoided. we need some way to define multiple response types per handler
			// and struct fields is the only way to achieve that.
			responseField := reflect.ValueOf(r.Response).FieldByIndex(responseType.fieldIndex).Interface()
			responseBytes, err = contentType.Encode(responseField)
			if err != nil {
				return RawResponse{}, err
			}
			r.Headers.Set(ContentTypeHeader, contentType.Mime())
		}
		bindResponseHeaders(ctx.Writer, r)
		ctx.Writer.WriteHeader(r.Status)
		_, err = ctx.Writer.Write(responseBytes)
		ctx.RawResponse.Status = r.Status
		ctx.RawResponse.Body = responseBytes
		ctx.RawResponse.Headers = r.Headers
		return *ctx.RawResponse, err
	}
}

// bindResponseHeaders copies Response.Headers to http.ResponseWriter.Header
func bindResponseHeaders[R any](writer http.ResponseWriter, r Response[R]) {
	for header, values := range r.Headers {
		for _, value := range values {
			writer.Header().Add(header, value)
		}
	}
}

// requestContentType extracts the ContentType implementation to use base on the "Content-Type" request header.
// If "Content-Type" request header is missing fallback to the provided default ContentType.
func requestContentType(r *http.Request, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	mimeType := r.Header.Get(ContentTypeHeader)
	if mimeType == "*/*" || mimeType == "" {
		return defaultContentType, nil
	}
	if contentType, found := supportedTypes[mimeType]; found {
		return contentType, nil
	}
	return nil, fmt.Errorf("unsupported mime type %q in %s header", mimeType, ContentTypeHeader)
}

// responseContentType extracts the ContentType implementation to use base on the "Accept" request header.
// If the "Accept" header is missing fallback to "Content-Type" request header.
// If "Content-Type" request header is also missing fallback to the provided default ContentType.
func responseContentType(r *http.Request, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	mimeTypes := []string{r.Header.Get(AcceptHeader), r.Header.Get(ContentTypeHeader)}
	for _, mimeType := range mimeTypes {
		if mimeType == "*/*" || mimeType == "" {
			return defaultContentType, nil
		}
		if contentTypes, found := supportedTypes[mimeType]; found {
			return contentTypes, nil
		}
	}
	return defaultContentType, fmt.Errorf("unsupported mime type %q in %s header", r.Header.Get(AcceptHeader), AcceptHeader)
}
