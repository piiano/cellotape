package utils

//type MultiError interface {
//	error
//	Errors() []error
//	Warnings() []error
//}
//
//type ErrorsCollector interface {
//	AddIfNotNil(LogLevel, ...error) bool
//	AddErrorsIfNotNil(...error) bool
//	AddWarningsIfNotNil(...error) bool
//	Errors() []error
//	Warnings() []error
//	ErrorOrNil() MultiError
//}
//
//func NewErrorsCollector() ErrorsCollector {
//	return &errorsCollector{
//		errors:   make([]error, 0),
//		warnings: make([]error, 0),
//	}
//}
//
//type errorsCollector struct {
//	errors   []error
//	warnings []error
//}
//
//func (e *errorsCollector) AddErrorsIfNotNil(errors ...error) bool {
//	return e.AddIfNotNil(Error, errors...)
//}
//func (e *errorsCollector) AddWarningsIfNotNil(errors ...error) bool {
//	return e.AddIfNotNil(Warn, errors...)
//}
//
//// AddIfNotNil check if errors are non nil values.
//// For non-nil errors add to the collector.
//// Return boolean value indicating an error was added to the collector or not.
//func (e *errorsCollector) AddIfNotNil(action LogLevel, errors ...error) bool {
//	if action == Ignore {
//		return false
//	}
//	added := false
//	for _, err := range errors {
//		if err == nil {
//			continue
//		}
//		if multiError, ok := err.(MultiError); ok {
//			e.AddIfNotNil(Warn, multiError.Warnings()...)
//			if e.AddIfNotNil(action, multiError.Errors()...) {
//				added = true
//			}
//			continue
//		}
//		switch action {
//		case Warn:
//			e.warnings = append(e.warnings, err)
//		case Error:
//			e.errors = append(e.errors, err)
//			added = true
//		}
//	}
//	return added
//}
//
//// Errors return all errors collected by the collector
//func (e *errorsCollector) Errors() []error {
//	return e.errors
//}
//
//// Errors return all errors collected by the collector
//func (e *errorsCollector) Warnings() []error {
//	return e.warnings
//}
//
//// ErrorOrNil if the collector has no errors return nil.
//// Else return an error that describes all errors collected in the collector.
//func (e *errorsCollector) ErrorOrNil() MultiError {
//	if len(e.errors) == 0 {
//		return nil
//	}
//	return e
//}
//
//// ErrorOrNil if the collector has no errors return nil.
//// Else return an error that describes all errors collected in the collector.
//func (e *errorsCollector) Error() string {
//	if len(e.errors) == 0 {
//		return ""
//	}
//	if len(e.errors) == 1 {
//		return fmt.Sprintf("1 error occurred:\n\t* %s\n", e.errors[0])
//	}
//	lines := make([]string, len(e.errors))
//	lines = append(lines, fmt.Sprintf("%d errors occurred:", len(e.errors)))
//	for _, err := range e.errors {
//		lines = append(lines, fmt.Sprintf("\t* %s", err))
//	}
//	lines = append(lines, "")
//	return strings.Join(lines, "\n")
//}
