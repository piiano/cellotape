package router

import "reflect"

//type Response[R any] interface {
//	error
//	Status()
//	Value() R
//}

type OK struct{ message string }
type BadRequest struct{ reason string }
type NotFound struct{ resource string }

type Response struct {
	OK         `status:"200"`
	BadRequest `status:"400"`
	NotFound   `status:"404"`
}

var name Response = Response{
	BadRequest: BadRequest{},
}

type responseTypes struct {
	responses map[int]reflect.Type
}
