package router

import (
	"fmt"
	"reflect"
)

type HttpResponse interface {
	getContentType() string
	getStatus() int
	getBodyType() (reflect.Type, error)
	bytes(contentTypes map[string]ContentType) ([]byte, error)
}

type ResponseType reflect.Type

func NewHttpResponse[R any](status int, contentType string) HttpResponse {
	return httpResponse[R]{
		status:      status,
		contentType: contentType,
	}
}

type httpResponse[R any] struct {
	status      int
	contentType string
	body        R
}

func (r httpResponse[R]) getContentType() string {
	return r.contentType
}

func (r httpResponse[R]) getStatus() int {
	return r.status
}

func (r httpResponse[R]) bytes(contentTypes map[string]ContentType) ([]byte, error) {
	contentType, found := contentTypes[r.contentType]
	if !found {
		return nil, fmt.Errorf("missing definition for content type %q", r.contentType)
	}
	return contentType.Marshal(r.body)
}

func (r httpResponse[R]) getBodyType() (reflect.Type, error) {
	var body R
	return reflect.TypeOf(body), nil
}

/*



 */

//func Resp[]() O {
//
//}
