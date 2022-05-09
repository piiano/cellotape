package router

import (
	"fmt"
	"reflect"
)

func failedValidatingTheRouterWithTheSpec(warnings int, errors int) string {
	return fmt.Sprintf("failed validating the router with the spec (%d warnings, %d errors)", warnings, errors)
}

func multipleHandlersFoundForOperationId(id string) string {
	return fmt.Sprintf("multiple handlers found for operation id %q", id)
}

func missingHandlerForOperationId(id string) string {
	return fmt.Sprintf("missing handler for operation id %q", id)
}

func notImplementedSpecOperations(count int) string {
	return fmt.Sprintf("%d of the spec operations are missing handlers", count)
}

func missingContentTypeImplementation(contentType string) string {
	return fmt.Sprintf("content type %s is declared in the spec but has no implementation in the router", contentType)
}

func paramDefinedByHandlerButMissingInSpec(in string, name string, paramsType reflect.Type, operationId string) string {
	return fmt.Sprintf("%s param %q is defined by type %s for operation %s but is not defined in the spec for that operation", in, name, paramsType, operationId)
}

func handlerDefinesRequestBodyWhenNoRequestBodyInSpec(operationID string) string {
	return fmt.Sprintf("handler defines a request body for operation %s while in the spec there is no request body for this operation", operationID)
}

func handlerForNonExistingSpecOperation(operationId string, position sourcePosition) string {
	return fmt.Sprintf("handler received for non exising operation id %q is spec - %s", operationId, position)
}

func invalidStatusInSpecResponses(statusStr string, operationId string) string {
	return fmt.Sprintf("spec declares an invalid status %s on operation %s", statusStr, operationId)
}

func handlerDefinesResponseThatIsMissingInSpec(status int, operationId string) string {
	return fmt.Sprintf("response %d is declared on operation %s but is not declared in the spec", status, operationId)
}

func unimplementedResponsesForOperation(unimplementedResponses int, operationId string) string {
	return fmt.Sprintf("%d responses exist on the spec for operation %s but not declared on any handler", unimplementedResponses, operationId)
}

func handlerDefinesResponseThatIsMissingInTheSpec(status int, operationId string) string {
	return fmt.Sprintf("response %d is declared on an handler for operation %s but is not part of the spec", status, operationId)
}
