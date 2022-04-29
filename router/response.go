package router

import (
	"fmt"
	"math/bits"
	"reflect"
	"strconv"
)

type ResponseType reflect.Type

type response struct {
	status       int
	responseType reflect.Type
	fieldIndex   []int
	isNilType    bool
}

type Send[T any] func(status int, response T)

func extractResponses(t reflect.Type) (map[int]response, error) {
	responseTypesMap := make(map[int]response, 0)
	if t.Kind() != reflect.Struct {
		return responseTypesMap, fmt.Errorf("responses type %s is not a struct type", t)
	}
	for _, field := range reflect.VisibleFields(t) {
		// only look at direct fields and not fields of embedded structs
		if len(field.Index) != 1 {
			continue
		}
		// each direct field of the responses struct need to have a status tag
		statusTag, ok := field.Tag.Lookup("status")
		if !ok {
			return responseTypesMap, fmt.Errorf("field %s of responses type %s is missing a status tag", field.Name, t.String())
		}
		status, err := parseStatus(statusTag)
		if err != nil {
			return responseTypesMap, err
		}
		//return responseTypesMap, fmt.Errorf("invalid status tag value %q for field %s of responses type %s", statusTag, field.Name, t.String())

		// each field represent a possible response
		responseTypesMap[int(status)] = response{
			status:       int(status),
			fieldIndex:   field.Index,
			responseType: field.Type,
			isNilType:    field.Type == nilType,
		}
	}
	return responseTypesMap, nil
}

func parseStatus(statusString string) (int, error) {
	status, err := strconv.ParseInt(statusString, 10, bits.UintSize)
	if err != nil || status < 100 || status >= 600 {
		return 0, fmt.Errorf("invalid status code %q", statusString)
	}
	return int(status), nil
}
