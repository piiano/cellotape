package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"reflect"
)

type handlerBinders[B, P, Q, R any] struct {
	requestBodyBinder        func(*gin.Context, *B) error
	pathParamsBinder         func(*gin.Context, *P) error
	queryParamsBinder        func(*gin.Context, *Q) error
	successfulResponseBinder func(*gin.Context, *R)
	errorResponseBinder      func(*gin.Context, *error)
}

func bindersFactory[B, P, Q, R any](oa OpenAPI, fn operationFunc[B, P, Q, R]) handlerBinders[B, P, Q, R] {
	types, _ := fn.types()
	var binders = handlerBinders[B, P, Q, R]{
		requestBodyBinder:        requestBodyBinderFactory[B](types.RequestBody, oa.getContentTypes()),
		pathParamsBinder:         pathBinderFactory[P](types.PathParams),
		queryParamsBinder:        queryBinderFactory[Q](types.QueryParams),
		successfulResponseBinder: responseBinderFactory[R](200, types.ResponseType, oa.getContentTypes()),
		errorResponseBinder:      responseBinderFactory[error](500, reflect.TypeOf(errors.New("")), oa.getContentTypes()),
	}
	return binders
}

func responseBinderFactory[R any](status int, responseType reflect.Type, contentTypes ContentTypes) func(*gin.Context, *R) {
	if responseType == nilType {
		return func(c *gin.Context, response *R) {
			c.Status(status)
		}
	}
	return func(c *gin.Context, response *R) {
		if status >= 400 {
			fmt.Printf("[ERROR] %+v\n", *response)
		}
		contentType, err := responseContentType(c, contentTypes, JsonContentType{})
		if err != nil {
			fmt.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
		}
		responseBytes, err := contentType.Marshal(response)
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err)
			fmt.Printf("[ERROR] failed serializing response %+v for mime type %s\n", *response, contentType.Mime())
			c.Data(500, contentType.Mime(), []byte(`{ "error": "failed serializing error response" }`))
			return
		}
		c.Data(status, contentType.Mime(), responseBytes)
	}
}

func requestBodyBinderFactory[B any](requestBodyType reflect.Type, contentTypes ContentTypes) func(*gin.Context, *B) error {
	if requestBodyType == nilType {
		return func(c *gin.Context, body *B) error {
			if c.Request.ContentLength != 0 {
				return errors.New("expected request with no body payload")
			}
			return nil
		}
	}
	return func(c *gin.Context, body *B) error {
		contentType, err := requestContentType(c, contentTypes, JsonContentType{})
		if err != nil {
			return err
		}
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		if err = contentType.Unmarshal(bodyBytes, body); err != nil {
			return err
		}
		return nil
	}
}

func pathBinderFactory[P any](pathParamsType reflect.Type) func(*gin.Context, *P) error {
	if pathParamsType == nilType {
		return func(c *gin.Context, body *P) error {
			if len(c.Params) > 0 {
				return fmt.Errorf("expected no path params but received %d", len(c.Params))
			}
			return nil
		}
	}
	return func(c *gin.Context, pathParams *P) error {
		if err := c.BindUri(pathParams); err != nil {
			return err
		}
		return nil
	}
}

func queryBinderFactory[Q any](queryParamsType reflect.Type) func(*gin.Context, *Q) error {
	if queryParamsType == nilType {
		// do nothing if there are query params in the request when no params expected
		return func(*gin.Context, *Q) error { return nil }
	}
	return func(c *gin.Context, queryParams *Q) error {
		if err := c.BindQuery(queryParams); err != nil {
			return err
		}
		return nil
	}
}

func requestContentType(c *gin.Context, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	mimeType := c.ContentType()
	if mimeType == "*/*" {
		return defaultContentType, nil
	}
	if contentType, found := supportedTypes[mimeType]; found {
		return contentType, nil
	}
	return nil, fmt.Errorf("unsupported mime type %q in Content-Type header", mimeType)
}

func responseContentType(c *gin.Context, supportedTypes ContentTypes, defaultContentType ContentType) (ContentType, error) {
	mimeTypes := []string{c.GetHeader("Accept"), c.ContentType()}
	for _, mimeType := range mimeTypes {
		if mimeType == "*/*" {
			return defaultContentType, nil
		}
		if contentTypes, found := supportedTypes[mimeType]; found {
			return contentTypes, nil
		}
	}
	return defaultContentType, fmt.Errorf("unsupported mime type %q in Accept header", c.GetHeader("Accept"))
}
