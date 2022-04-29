package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"reflect"
)

type handlerBinders[B, P, Q, R any] struct {
	requestBodyBinder func(*gin.Context, *B) error
	pathParamsBinder  func(*gin.Context, *P) error
	queryParamsBinder func(*gin.Context, *Q) error
	sendFactory       func(*gin.Context) Send[R]
	sendError         func(*gin.Context, error)
}

type internalServerError struct {
	Error string `json:"error"`
}

func bindersFactory[B, P, Q, R any](oa OpenAPI, fn operationFunc[B, P, Q, R]) handlerBinders[B, P, Q, R] {
	types := fn.types()
	var binders = handlerBinders[B, P, Q, R]{
		requestBodyBinder: requestBodyBinderFactory[B](types.requestBody, oa.getContentTypes()),
		pathParamsBinder:  pathBinderFactory[P](types.pathParams),
		queryParamsBinder: queryBinderFactory[Q](types.queryParams),
		sendFactory:       responseBinderFactory[R](types.responsesType, oa.getContentTypes()),
		sendError:         sendErrorFactory(oa.getContentTypes()),
	}
	return binders
}

func sendErrorFactory(contentTypes ContentTypes) func(*gin.Context, error) {
	return func(c *gin.Context, err error) {
		contentType, err := responseContentType(c, contentTypes, JsonContentType{})
		if err != nil {
			fmt.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
		}
		errResponse := internalServerError{Error: err.Error()}
		responseBytes, err := contentType.Marshal(errResponse)
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err)
			fmt.Printf("[ERROR] failed serializing response %+v for mime type %s\n", errResponse, contentType.Mime())
			c.Data(500, contentType.Mime(), []byte(`{ "error": "failed serializing error response" }`))
			return
		}
		c.Data(500, contentType.Mime(), responseBytes)
	}
}

func responseBinderFactory[R any](responsesType reflect.Type, contentTypes ContentTypes) func(*gin.Context) Send[R] {
	responseTypesMap, err := extractResponses(responsesType)
	if err != nil {

	}
	return func(c *gin.Context) Send[R] {
		return func(status int, responses R) {
			responseType, _ := responseTypesMap[status]
			if responseType.isNilType {
				c.Status(status)
				return
			}
			response := reflect.ValueOf(responses).FieldByIndex(responseType.fieldIndex).Interface()
			contentType, err := responseContentType(c, contentTypes, JsonContentType{})
			if err != nil {
				fmt.Printf("[WARNING] %s. fallback to %s\n", err, contentType.Mime())
			}
			responseBytes, err := contentType.Marshal(response)
			if err != nil {
				fmt.Printf("[ERROR] %s\n", err)
				fmt.Printf("[ERROR] failed serializing response %+v for mime type %s\n", response, contentType.Mime())
				c.Data(500, contentType.Mime(), []byte(`{ "error": "failed serializing error response" }`))
				return
			}
			c.Data(status, contentType.Mime(), responseBytes)
		}
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
