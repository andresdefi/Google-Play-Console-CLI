package exitcode

const (
	Success  = 0
	Error    = 1
	Usage    = 2
	Auth     = 3
	NotFound = 4
	Conflict = 5
	Config   = 6

	// HTTP 4xx range: 10 + (status - 400), except special cases above.
	// HTTP 5xx range: 60 + (status - 500).
	http4xxBase = 10
	http5xxBase = 60
)

// FromHTTPStatus maps an HTTP status code to a granular exit code.
// Special cases: 401/403->Auth, 404->NotFound, 409->Conflict.
// General 4xx maps to 10-59, 5xx maps to 60-99.
func FromHTTPStatus(status int) int {
	switch {
	case status >= 200 && status < 300:
		return Success
	case status == 401, status == 403:
		return Auth
	case status == 404:
		return NotFound
	case status == 409:
		return Conflict
	case status >= 400 && status < 500:
		code := http4xxBase + (status - 400)
		if code > 59 {
			code = 59
		}
		return code
	case status >= 500 && status < 600:
		code := http5xxBase + (status - 500)
		if code > 99 {
			code = 99
		}
		return code
	default:
		return Error
	}
}

// API is kept as an alias for backward compatibility in commands that
// don't need granular codes. Equivalent to NotFound.
const API = NotFound
