package types

type ErrorType int

const (
	// ErrorTypeValidation is returned when the request is invalid.
	ErrorTypeValidation ErrorType = iota

	// ErrorTypeSyntax is returned when the request has a syntax error.
	ErrorTypeSyntax

	// ErrorTypeInvalidCredentials is returned when the credentials are invalid.
	ErrorTypeInvalidCredentials

	// ErrorTypeUnauthorized is returned when the request is unauthorized.
	ErrorTypeUnauthorized

	// ErrorTypeNotFound is returned when the resource is not found.
	ErrorTypeNotFound

	// ErrorTypeConflict is returned when the request conflicts with the current state of the server.
	ErrorTypeConflict

	// ErrorTypeInternal is returned when the server has an internal error.
	ErrorTypeInternal

	// ErrorTypeServiceUnavailable is returned when the server is unavailable.
	ErrorTypeServiceUnavailable

	// ErrorTypeUnknown is returned when the server has an unknown error.
	ErrorTypeUnknown

	// ErrorBadRequest is returned when the request is bad.
	ErrorBadRequest

	// ErrorTypeAccountDeleted is returned when the account is deleted.
	ErrorTypeAccountDeleted

	// ErrorInvalidRequest is returned when the request is invalid.
	ErrorInvalidRequest
)

func (t ErrorType) String() string {
	return [...]string{
		"validation",
		"syntax",
		"invalid_credentials",
		"unauthorized",
		"not_found",
		"conflict",
		"internal",
		"service_unavailable",
		"unknown",
		"bad_request",
		"account_deleted",
		"invalid_request",
	}[t]
}

type Status int

const (
	// Success is returned when the request is successful.
	Success Status = iota

	// Pending is returned when the request is pending.
	Pending

	// Failure is returned when the request has failed.
	Failure
)

func (s Status) String() string {
	return [...]string{
		"success",
		"pending",
		"failure",
	}[s]
}

// ApiResponse represents an API response with status, code, message, error, and data.
type ApiResponse struct {
	Status  string      `json:"status,omitempty"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Type    string      `json:"type,omitempty"`
	Errors  []Error     `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason,omitempty"`
}
