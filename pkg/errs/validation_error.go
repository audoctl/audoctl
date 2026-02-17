package errs

// ValidationError struct uses for validation error messages
type ValidationError struct {
	Key     string         `json:"key"`
	Field   string         `json:"field"`
	Message string         `json:"message"`
	Params  map[string]any `json:"params,omitempty"`
}

// NewValidationError initializes validation error
func NewValidationError(key, field, message string, params map[string]any) ValidationError {
	if params == nil {
		params = make(map[string]any)
	}
	return ValidationError{
		Key:     key,
		Field:   field,
		Params:  params,
		Message: message,
	}
}

func (e *ValidationError) Error() string {
	return e.Message
}
