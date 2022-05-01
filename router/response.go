package router

import (
	"fmt"
	"math/bits"
	"reflect"
	"strconv"
)

type httpResponse struct {
	status       int
	responseType reflect.Type
	fieldIndex   []int
	isNilType    bool
}

func extractResponses(t reflect.Type) (map[int]httpResponse, error) {
	responseTypesMap := make(map[int]httpResponse, 0)
	if t.Kind() != reflect.Struct {
		return responseTypesMap, fmt.Errorf("responses type %s is not a struct type", t)
	}
	for _, field := range reflect.VisibleFields(t) {
		// only look at direct exported fields and not fields of embedded structs
		if len(field.Index) != 1 || !field.IsExported() {
			continue
		}

		// each direct field of the responses' struct need to have a status tag
		statusTag, ok := field.Tag.Lookup("status")
		if !ok {
			return responseTypesMap, fmt.Errorf("field %s of responses type %s is missing a status tag", field.Name, t.String())
		}
		status, err := parseStatus(statusTag)
		if err != nil {
			return responseTypesMap, err
		}
		// each field represent a possible httpResponse
		responseTypesMap[status] = httpResponse{
			status:       status,
			fieldIndex:   field.Index,
			responseType: field.Type,
			isNilType:    field.Type == nilType,
		}
	}
	return responseTypesMap, nil
}

// parse a string representing an HTTP status code or error if it is not a valid code between 100 and 600
func parseStatus(statusString string) (int, error) {
	status, err := strconv.ParseInt(statusString, 10, bits.UintSize)
	if err != nil || status < 100 || status >= 600 {
		return 0, fmt.Errorf("invalid status code %q", statusString)
	}
	return int(status), nil
}
