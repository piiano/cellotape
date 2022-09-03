package router

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	assert.Equal(t,
		"failed validating the router with the spec (5 warnings, 3 errors)",
		failedValidatingTheRouterWithTheSpec(5, 3))

	assert.Equal(t,
		`multiple handlers found for operation "foo"`,
		multipleHandlersFoundForOperationId("foo"))

	assert.Equal(t,
		`missing handler for operation "foo"`,
		missingHandlerForOperationId("foo"))

	assert.Equal(t,
		"3 of the spec operations are missing handlers",
		notImplementedSpecOperations(3))

	assert.Equal(t,
		"content type text/plain is declared in the spec but has no implementation in the router",
		missingContentTypeImplementation("text/plain"))

	assert.Equal(t,
		`query param "bar" is defined by type string for operation foo but is not defined in the spec for that operation`,
		paramDefinedByHandlerButMissingInSpec("query", "bar", reflect.TypeOf(""), "foo"))

	assert.Equal(t,
		"handler defines a request body for operation foo while in the spec there is no request body for this operation",
		handlerDefinesRequestBodyWhenNoRequestBodyInSpec("foo"))

	assert.Equal(t,
		`handler received for non exising operation "foo" is spec - file.go:10`,
		handlerForNonExistingSpecOperation("foo", sourcePosition{
			ok:   true,
			file: "file.go",
			line: 10,
		}))

	assert.Equal(t,
		"spec declares an invalid status 40x on operation foo",
		invalidStatusInSpecResponses("40x", "foo"))

	assert.Equal(t,
		"200 response exist on the spec for operation foo but not declared on any handler",
		unimplementedResponseForOperation(200, "foo"))

	assert.Equal(t,
		"response 200 is declared on a handler for operation foo but is not part of the spec",
		handlerDefinesResponseThatIsMissingInTheSpec(200, "foo"))

	assert.Equal(t,
		`request body schema of operation "foo" is incompatible with handler request body type string`,
		incompatibleRequestBodyType("foo", reflect.TypeOf("")))

	assert.Equal(t,
		`schema of query param "bar" of operation "foo" is incompatible with handler request param type string of field "Bar"`,
		incompatibleParamType("foo", "query", "bar", "Bar", reflect.TypeOf("")))

	assert.Equal(t,
		`200 response schema of operation "foo" is incompatible with handler 200 response type string`,
		incompatibleResponseType("foo", 200, reflect.TypeOf("")))

	assert.Equal(t,
		"200 responses exist on the spec for operation foo but not declared on any handler",
		unimplementedResponsesForOperation(200, "foo"))

	assert.Equal(t,
		"response 200 is declared on operation foo but is not declared in the spec",
		handlerDefinesResponseThatIsMissingInSpec(200, "foo"))

	assert.Equal(t,
		`query param "bar" exists on the spec for operation "foo" but not declared on any handler`,
		paramMissingImplementationInChain("query", "bar", "foo"))

}
