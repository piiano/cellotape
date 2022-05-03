package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"reflect"
)

type handlerBinders[B, P, Q, R any] struct {
	requestBinder  func(HandlerContext) (Request[B, P, Q], error)
	responseBinder func(HandlerContext, Response[R]) (Response[any], error)
}

// produce set of binder functions that can be called at runtime to handle each request
func bindersFactory[B, P, Q, R any](oa openapi, fn operationFunc[B, P, Q, R]) handlerBinders[B, P, Q, R] {
	return handlerBinders[B, P, Q, R]{
		requestBinder:  requestBinderFactory[B, P, Q](oa, fn.requestTypes()),
		responseBinder: responseBinderFactory[R](fn.responseTypes(), oa.contentTypes),
	}
}

// produce the binder function that can be called at runtime to create the request object for the handler
func requestBinderFactory[B, P, Q any](oa openapi, types handlerRequestTypes) func(HandlerContext) (Request[B, P, Q], error) {
	requestBodyBinder := requestBodyBinderFactory[B](types.requestBody, oa.contentTypes)
	pathParamsBinder := pathBinderFactory[P](types.pathParams)
	queryParamsBinder := queryBinderFactory[Q](types.queryParams)

	// this is what actually build the request object at runtime for the handler
	return func(c HandlerContext) (Request[B, P, Q], error) {
		var request = Request[B, P, Q]{Context: c.Request.Context(), Headers: c.Request.Header}
		if err := requestBodyBinder(c.Request, &request.Body); err != nil {
			return request, err
		}
		if err := pathParamsBinder(c.Params, &request.PathParams); err != nil {
			return request, err
		}
		if err := queryParamsBinder(c.Request, &request.QueryParams); err != nil {
			return request, err
		}
		return request, nil
	}
}

// produce the request body binder that can be used in runtime
func requestBodyBinderFactory[B any](requestBodyType reflect.Type, contentTypes ContentTypes) func(*http.Request, *B) error {
	if requestBodyType == nilType {
		return func(r *http.Request, body *B) error {
			if r.ContentLength != 0 {
				return errors.New("expected request with no body payload")
			}
			return nil
		}
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

// produce the path params binder that can be used in runtime
func pathBinderFactory[P any](pathParamsType reflect.Type) func(httprouter.Params, *P) error {
	if pathParamsType == nilType {
		return func(params httprouter.Params, body *P) error {
			if len(params) > 0 {
				return fmt.Errorf("expected no path params but received %d", len(params))
			}
			return nil
		}
	}
	return func(params httprouter.Params, pathParams *P) error {
		m := make(map[string][]string)
		for _, v := range params {
			m[v.Key] = []string{v.Value}
		}
		return binding.Uri.BindUri(m, pathParams)
	}
}

// produce the query params binder that can be used in runtime
func queryBinderFactory[Q any](queryParamsType reflect.Type) func(*http.Request, *Q) error {
	if queryParamsType == nilType {
		// do nothing if there are query params in the request when no params expected
		return func(*http.Request, *Q) error { return nil }
	}
	return func(r *http.Request, queryParams *Q) error {
		if err := binding.Query.Bind(r, queryParams); err != nil {
			return err
		}
		return nil
	}
}
func responseBinderFactory[R any](responseTypes handlerResponseTypes, contentTypes ContentTypes) func(HandlerContext, Response[R]) (Response[any], error) {
	return func(c HandlerContext, r Response[R]) (Response[any], error) {
		response := Response[any]{
			Status:  r.Status,
			Headers: r.Headers,
			Bytes:   r.Bytes,
			bound:   r.bound,
		}
		if r.bound {
			return response, nil
		}
		contentType, err := responseContentType(c.Request, contentTypes, JsonContentType{})
		if err != nil {
			log.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
		}
		responseType, exist := responseTypes.declaredResponses[r.Status]
		if !exist {
			return response, fmt.Errorf("status %d is not part of the possible operation responses", r.Status)
		}
		if responseType.isNilType {
			response.bound = true
			return response, nil
		}
		responseField := reflect.ValueOf(r.Response).FieldByIndex(responseType.fieldIndex).Interface()
		responseBytes, err := contentType.Encode(responseField)
		if err != nil {
			return response, err
		}
		response.Headers.Set("Content-Type", contentType.Mime())
		response.Bytes = responseBytes
		response.bound = true
		return response, nil
	}
}

func requestContentType(r *http.Request, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	mimeType := r.Header.Get("Content-Type")
	if mimeType == "*/*" {
		return defaultContentType, nil
	}
	if contentType, found := supportedTypes[mimeType]; found {
		return contentType, nil
	}
	return nil, fmt.Errorf("unsupported mime type %q in Content-Type header", mimeType)
}

func responseContentType(r *http.Request, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	mimeTypes := []string{r.Header.Get("Accept"), r.Header.Get("Content-Type")}
	for _, mimeType := range mimeTypes {
		if mimeType == "*/*" {
			return defaultContentType, nil
		}
		if contentTypes, found := supportedTypes[mimeType]; found {
			return contentTypes, nil
		}
	}
	return defaultContentType, fmt.Errorf("unsupported mime type %q in Accept header", r.Header.Get("Accept"))
}
