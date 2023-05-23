package router

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

func TestStringFormatValidation(t *testing.T) {
	uuidSchema := openapi3.NewStringSchema().WithFormat("uuid")

	// Register additional validations
	registerAdditionalOpenAPIFormatValidations()

	// Make sure we can validate uuid after we register them.
	err := uuidSchema.VisitJSON("4083aee7-80aa-40be-8e36-e159484ff431", openapi3.EnableFormatValidation())
	require.NoError(t, err)

	// Make sure invalid uuid is now failing after we register the validations.
	err = uuidSchema.VisitJSON("not-a-uuid", openapi3.EnableFormatValidation())
	require.ErrorIs(t, err, ErrInvalidUUID)
}
