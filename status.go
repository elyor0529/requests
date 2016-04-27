package requests

type statusCode int

// Successful
const (
	Ok statusCode = 200 + iota
	Created
	Accepted
	_
	NoContent
)

// Redirection
const (
	MovedPermanently statusCode = 301 + iota
	Found
	SeeOther
	NotModified
	UseProxy
)

// Client Error
const (
	BadRequest statusCode = 400 + iota
	Unauthorized
	PaymentRequired
	Forbidden
	NotFound
	MethodNotAllowed
	NotAcceptable
	_
	RequestTimeout
	Conflict
)

// Server Error
const (
	InternalServerError statusCode = 500 + iota
	NotImplemented
	BadGateway
	ServiceUnavailable
)
