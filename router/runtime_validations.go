package router

import (
	"errors"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
)

var (
	// registerOpenAPIFormatValidationsOnce is used to make sure we register the validations only once.
	registerOpenAPIFormatValidationsOnce = &sync.Once{}

	// ErrInvalidUUID is returned when the uuid format validation fails.
	ErrInvalidUUID = errors.New("invalid uuid")
)

// registerAdditionalOpenAPIFormatValidations registers additional validations that are not supported by the kin-openapi.
// kin-openapi supports out of the box the following formats: `date-time`, `date` and `byte`.
func registerAdditionalOpenAPIFormatValidations() {

	// Make sure we register the validations only once (even if we have multiple openapi instances)
	// This is because the kin-openapi library uses a global map to store the validations and we don't want to have a
	// race condition when initializing multiple routers.
	registerOpenAPIFormatValidationsOnce.Do(func() {

		// Register uuid format validation
		openapi3.DefineStringFormatCallback("uuid", func(value string) error {
			if _, err := uuid.Parse(value); err != nil {
				return ErrInvalidUUID
			}

			return nil
		})
	})
}
