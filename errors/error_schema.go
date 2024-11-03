package e


type ErrorDetail struct {
	Field string `json:"field,omitempty"` // Field causing the error
	Issue string `json:"issue,omitempty"` // Description of the issue with the field
}

// ErrorResponse struct for API error responses
type ErrorResponse struct {
	Code          string        `json:"code"`                    // Error code
	Message       string        `json:"message"`                 // Human-readable error message
	Error			string		`json:"error"`                 	// Stack Trace
	Details       []ErrorDetail `json:"details,omitempty"`       // List of field-specific errors
}