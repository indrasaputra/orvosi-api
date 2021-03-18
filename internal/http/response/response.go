package response

// Success represents success response.
type Success struct {
	// Data represents any primary data that is visible in response.
	Data interface{} `json:"data"`
	// Meta represents auxiliary data that is visible in response.
	Meta interface{} `json:"meta"`
}

// Error represents error response.
type Error struct {
	// Errors represents list of errors that are visible in response.
	Errors []error `json:"errors"`
	// Meta represents auxiliary data that is visible in response.
	Meta interface{} `json:"meta"`
}

// EmptyMeta represents an empty struct.
type EmptyMeta struct{}

// NewSuccess creates an instance of Success response.
func NewSuccess(data, meta interface{}) *Success {
	return &Success{
		Data: data,
		Meta: meta,
	}
}

// NewError creates an instance of Error response.
func NewError(errors ...error) *Error {
	return &Error{
		Errors: errors,
		Meta:   nil,
	}
}
