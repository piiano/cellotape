package router

import (
	"fmt"
	"math/bits"
	"reflect"
	"strconv"

	"github.com/piiano/cellotape/router/utils"
)

const statusTag = "status"

// extractResponses only extracts gracefully a map of responses declared in a type.
// ignore responses formatted badly or return an empty map of the entire type isn't representing a valid response type
func extractResponses(t reflect.Type) handlerResponses {
	responseTypesMap := make(handlerResponses, 0)
	if t == nil || t.Kind() != reflect.Struct {
		return responseTypesMap
	}
	for _, field := range reflect.VisibleFields(t) {
		//TODO: handle invalid fields

		// only look at direct exported fields and not fields of embedded structs
		if len(field.Index) != 1 || !field.IsExported() {
			continue
		}
		if field.Anonymous {
			for status, response := range extractResponses(field.Type) {
				response.fieldIndex = append(field.Index, response.fieldIndex...)
				responseTypesMap[status] = response
			}
			continue
		}
		// each direct field of the responses' struct need to have a status tag
		statusTagValue, ok := field.Tag.Lookup(statusTag)
		if !ok {
			continue
		}
		status, err := parseStatus(statusTagValue)
		if err != nil {
			continue
		}
		// each field represent a possible httpResponse
		responseTypesMap[status] = httpResponse{
			status:       status,
			fieldIndex:   field.Index,
			responseType: field.Type,
			isNilType:    field.Type == utils.NilType,
		}
	}
	return responseTypesMap
}

// parse a string representing an HTTP status code or error if it is not a valid code between 100 and 600
func parseStatus(statusString string) (int, error) {
	status, err := strconv.ParseInt(statusString, 10, bits.UintSize)
	if err != nil || status < 100 || status >= 600 {
		return 0, fmt.Errorf("invalid status code %q", statusString)
	}
	return int(status), nil
}
