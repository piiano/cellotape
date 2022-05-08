package utils

//import (
//	"errors"
//	"strings"
//	"testing"
//)
//
//func TestNewErrorsCollector(t *testing.T) {
//	errs := NewErrorsCollector()
//	if errs == nil {
//		t.Error("expect NewErrorsCollector() to return non nil value")
//	}
//	_, ok := errs.(ErrorsCollector)
//	if !ok {
//		t.Error("expect NewErrorsCollector() to return ErrorsCollector")
//	}
//}
//
//func TestAddIfNotNilWithNilValues(t *testing.T) {
//	errs := NewErrorsCollector()
//	errs.AddIfNotNil(nil)
//	errs.AddIfNotNil(nil, nil, nil)
//	if len(errs.Errors()) != 0 {
//		t.Error("nil errors should not be added to the collector")
//	}
//	if errs.ErrorOrNil() != nil {
//		t.Error("nil errors should not be added to the collector")
//	}
//}
//
//func TestAddIfNotNilWithMultiError(t *testing.T) {
//	errs := NewErrorsCollector()
//	var emptyMultiErr MultiError = &errorsCollector{
//		errors: []error{},
//	}
//	errs.AddIfNotNil(emptyMultiErr)
//	if len(errs.Errors()) != 0 {
//		t.Error("empty multi error should not be added")
//	}
//
//	errorMessage1 := "error message 1"
//	errorMessage2 := "error message 2"
//	errorMessage3 := "error message 3"
//	var multiErr MultiError = &errorsCollector{
//		errors: []error{
//			errors.New(errorMessage1),
//			errors.New(errorMessage2),
//			errors.New(errorMessage3),
//		},
//	}
//	errs.AddIfNotNil(multiErr)
//	if len(errs.Errors()) != 3 {
//		t.Error("multi error added should be flattened error")
//	}
//	errMessage := errs.ErrorOrNil().Error()
//	if !strings.Contains(errMessage, errorMessage1) ||
//		!strings.Contains(errMessage, errorMessage2) ||
//		!strings.Contains(errMessage, errorMessage3) {
//		t.Error("multi error added should preserve internal error messages")
//	}
//}
//
//func TestAddIfNotNilWithErrorValue(t *testing.T) {
//	errs := NewErrorsCollector()
//	errorMessage := "error message"
//	dummyErr := errors.New(errorMessage)
//	errs.AddIfNotNil(nil, dummyErr, nil)
//	if len(errs.Errors()) != 1 {
//		t.Error("the collector should have a single error in it")
//	}
//	if errs.Errors()[0] != dummyErr {
//		t.Error("the collector should contain the added error")
//	}
//	multiErr := errs.ErrorOrNil()
//	if multiErr == nil {
//		t.Error("expected to have return an error")
//	}
//	if !strings.Contains(multiErr.Error(), errorMessage) {
//		t.Errorf("expected the final multi error to include the added error message %q", errorMessage)
//	}
//	if len(multiErr.Errors()) != len(errs.Errors()) {
//		t.Error("expect the returned multiError to have the same errors as in the collector")
//	}
//	for i, err := range errs.Errors() {
//		if multiErr.Errors()[i] != err {
//			t.Error("expect the returned multiError to have the same errors as in the collector")
//		}
//	}
//}
//
//func TestAddIfNotNilWithMultipleErrorValues(t *testing.T) {
//	errs := NewErrorsCollector()
//	errorMessage1 := "error message 1"
//	errorMessage2 := "error message 2"
//	dummyErr1 := errors.New(errorMessage1)
//	dummyErr2 := errors.New(errorMessage2)
//	errs.AddIfNotNil(dummyErr1)
//	errs.AddIfNotNil(nil)
//	errs.AddIfNotNil(nil, dummyErr2, nil)
//	if len(errs.Errors()) != 2 {
//		t.Error("the collector should have a both errors in it")
//	}
//	if errs.Errors()[0] != dummyErr1 {
//		t.Error("the collector should contain the first error and preserve order")
//	}
//	if errs.Errors()[1] != dummyErr2 {
//		t.Error("the collector should contain the second error and preserve order")
//	}
//	multiErr := errs.ErrorOrNil()
//	if multiErr == nil {
//		t.Error("expected to have return an error")
//	}
//	if !strings.Contains(multiErr.Error(), errorMessage1) {
//		t.Errorf("expected the final multi error to include the added error message %q", errorMessage1)
//	}
//	if !strings.Contains(multiErr.Error(), errorMessage2) {
//		t.Errorf("expected the final multi error to include the added error message %q", errorMessage2)
//	}
//	if len(multiErr.Errors()) != len(errs.Errors()) {
//		t.Error("expect the returned multiError to have the same errors as in the collector")
//	}
//	for i, err := range errs.Errors() {
//		if multiErr.Errors()[i] != err {
//			t.Error("expect the returned multiError to have the same errors as in the collector")
//		}
//	}
//}
//
//func TestInternalErrorImplementationForNoErrors(t *testing.T) {
//	errs := &errorsCollector{
//		errors: make([]error, 0),
//	}
//	if errs.Error() != "" {
//		t.Error("expect internal collector implementation to return empty string when no errors")
//	}
//}
