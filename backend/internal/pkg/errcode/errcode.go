package errcode

// Business error codes — returned in Response.Code field.
// HTTP status codes are set separately at the transport layer.
const (
	// Success
	OK = 0

	// Generic client errors (4xx range)
	ErrBadRequest   = 40000
	ErrUnauthorized = 40001
	ErrForbidden    = 40003
	ErrNotFound     = 40004
	ErrConflict     = 40009

	// Validation
	ErrValidation = 42200

	// Business logic
	ErrCruiseHasCabins   = 42201
	ErrCompanyHasCruises = 42202
	ErrPasswordMismatch  = 42203

	// Server errors (5xx range)
	ErrInternal = 50000
)
