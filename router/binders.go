package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"reflect"
)

type handlerBinders[B, P, Q, R any] struct {
	requestBinder  func(*gin.Context) (Request[B, P, Q], error)
	responseBinder func(*gin.Context, int, R)
	errorBinder    func(*gin.Context, error)
}

type internalServerError struct {
	Error string `json:"error"`
}

// produce set of binder functions that can be called at runtime to handle each request
func bindersFactory[B, P, Q, R any](oa openapi, fn operationFunc[B, P, Q, R]) handlerBinders[B, P, Q, R] {
	types := fn.types()
	var binders = handlerBinders[B, P, Q, R]{
		requestBinder:  requestBinderFactory[B, P, Q](oa, types),
		responseBinder: responseBinderFactory[R](types.responsesType, oa.contentTypes),
		errorBinder:    errorBinderFactory(oa.contentTypes),
	}
	return binders
}

// produce the binder function that can be called at runtime to create the request object for the handler
func requestBinderFactory[B, P, Q any](oa openapi, types operationTypes) func(*gin.Context) (Request[B, P, Q], error) {
	requestBodyBinder := requestBodyBinderFactory[B](types.requestBody, oa.contentTypes)
	pathParamsBinder := pathBinderFactory[P](types.pathParams)
	queryParamsBinder := queryBinderFactory[Q](types.queryParams)

	// this is what actually build the request object at runtime for the handler
	return func(c *gin.Context) (Request[B, P, Q], error) {
		var request = Request[B, P, Q]{Context: c, Headers: c.Request.Header}
		if err := requestBodyBinder(c, &request.Body); err != nil {
			return request, err
		}
		if err := pathParamsBinder(c, &request.PathParams); err != nil {
			return request, err
		}
		if err := queryParamsBinder(c, &request.QueryParams); err != nil {
			return request, err
		}
		return request, nil
	}
}

// produce the request body binder that can be used in runtime
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
		if err = contentType.Decode(bodyBytes, body); err != nil {
			return err
		}
		return nil
	}
}

// produce the path params binder that can be used in runtime
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

// produce the query params binder that can be used in runtime
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
func responseBinderFactory[R any](responsesType reflect.Type, contentTypes ContentTypes) func(*gin.Context, int, R) {
	responseTypesMap, err := extractResponses(responsesType)
	if err != nil {

	}
	return func(c *gin.Context, status int, responses R) {
		responseType, _ := responseTypesMap[status]
		if responseType.isNilType {
			c.Status(status)
			return
		}
		response := reflect.ValueOf(responses).FieldByIndex(responseType.fieldIndex).Interface()
		contentType, err := responseContentType(c, contentTypes, JsonContentType{})
		if err != nil {
			log.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
		}
		responseBytes, err := contentType.Encode(response)
		if err != nil {
			log.Printf("[ERROR] %s\n", err)
			log.Printf("[ERROR] failed serializing httpResponse %+v for mime type %s\n", response, contentType.Mime())
			c.Data(500, contentType.Mime(), []byte(`{ "error": "failed serializing error httpResponse" }`))
			return
		}
		c.Data(status, contentType.Mime(), responseBytes)
	}
}

func errorBinderFactory(contentTypes ContentTypes) func(*gin.Context, error) {
	return func(c *gin.Context, errResponse error) {
		contentType, err := responseContentType(c, contentTypes, JsonContentType{})
		if err != nil {
			log.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
		}
		errHTTPResponse := internalServerError{Error: errResponse.Error()}
		responseBytes, err := contentType.Encode(errHTTPResponse)
		if err != nil {
			log.Printf("[ERROR] %s\n", err)
			log.Printf("[ERROR] failed serializing httpResponse %+v for mime type %s\n", errResponse, contentType.Mime())
			c.Data(500, contentType.Mime(), []byte(`{ "error": "failed serializing error httpResponse" }`))
			return
		}
		c.Data(500, contentType.Mime(), responseBytes)
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
